from fastapi import FastAPI, WebSocket, WebSocketDisconnect
from aiokafka import AIOKafkaProducer, AIOKafkaConsumer
from contextlib import asynccontextmanager
import json
import asyncio
import os

KAFKA_SERVER = os.getenv("KAFKA_SERVER", "localhost:9092")
KAFKA_TOPIC = "chat-messages"

class ConnectionManager:
    def __init__(self):
        self.active_connections: dict = {}

    async def connect(self, user_id: str, websocket: WebSocket):
        await websocket.accept()
        self.active_connections[user_id] = websocket

    def disconnect(self, user_id: str):
        if user_id in self.active_connections:
            del self.active_connections[user_id]

    async def send_personal_message(self, message: str, user_id: str):
        if user_id in self.active_connections:
            await self.active_connections[user_id].send_text(message)

manager = ConnectionManager()

async def consume_messages():
    consumer = AIOKafkaConsumer(
        KAFKA_TOPIC,
        bootstrap_servers=KAFKA_SERVER,
        group_id="sms-group"
    )
    await consumer.start()
    try:
        async for msg in consumer:
            data = json.loads(msg.value.decode("utf-8"))
            receiver = data.get("receiver")
            await manager.send_personal_message(
                json.dumps({"from": data['sender'], "msg": data['content']}), 
                receiver
            )
    finally:
        await consumer.stop()

@asynccontextmanager
async def lifespan(app: FastAPI):
    app.state.producer = AIOKafkaProducer(bootstrap_servers=KAFKA_SERVER)
    await app.state.producer.start()
    consumer_task = asyncio.create_task(consume_messages())
    yield
    await app.state.producer.stop()
    consumer_task.cancel()

app = FastAPI(lifespan=lifespan)

@app.websocket("/ws/{user_id}")
async def websocket_endpoint(websocket: WebSocket, user_id: str):
    await manager.connect(user_id, websocket)
    try:
        while True:
            data = await websocket.receive_text()
            msg_json = json.loads(data)
            message_to_kafka = {
                "sender": user_id,
                "receiver": msg_json["receiver"],
                "content": msg_json["content"]
            }
            await app.state.producer.send_and_wait(KAFKA_TOPIC, json.dumps(message_to_kafka).encode("utf-8"))
    except WebSocketDisconnect:
        manager.disconnect(user_id)
    except Exception as e:
        manager.disconnect(user_id)