import asyncio
from fastapi.concurrency import asynccontextmanager
from fastapi import Depends, FastAPI, HTTPException
from sqlalchemy import select
from python_api.db.session import SessionLocal
from python_api.services.rabbitmq import rabbitmq_listener
from python_api.entities.payment_entity import Payment
from python_api.models.payment import PaymentResponse, PaymentCreate
from sqlalchemy.ext.asyncio import AsyncSession

app = FastAPI()

async def get_db():
    async with SessionLocal() as session:
        yield session

@asynccontextmanager
async def lifespan(app: FastAPI):
    task = asyncio.create_task(rabbitmq_listener())
    yield
    task.cancel()
    try:
        await task
    except asyncio.CancelledError:
        pass

app = FastAPI(lifespan=lifespan)

@app.get("/ping")
async def ping():
    return "pong 🏓"

@app.post("/payments/", response_model=PaymentResponse)
async def create_payment(payment: PaymentCreate, db: AsyncSession = Depends(get_db)):
    new_payment = Payment(**payment.model_dump())
    db.add(new_payment)
    await db.commit()
    await db.refresh(new_payment)
    return PaymentResponse(
        id=new_payment.id, 
        user_id=new_payment.user_id, 
        amount=new_payment.amount, 
        description=new_payment.description
    )

@app.get("/payments/{payment_id}", response_model=PaymentResponse)
async def read_payment(payment_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Payment).where(Payment.id == payment_id))
    payment = result.scalars().first()
    if not payment:
        raise HTTPException(status_code=404, detail="Payment not found")
    return PaymentResponse(
        id=payment.id, 
        user_id=payment.user_id, 
        amount=payment.amount, 
        description=payment.description
    )
@app.put("/payments/{payment_id}", response_model=PaymentResponse)
async def update_payment(payment_id: int, payment_update: PaymentCreate, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Payment).where(Payment.id == payment_id))
    payment = result.scalars().first()
    if not payment:
        raise HTTPException(status_code=404, detail="Payment not found")
    for key, value in payment_update.model_dump().items():
        setattr(payment, key, value)
    await db.commit()
    await db.refresh(payment)
    return PaymentResponse(
        id=payment.id, 
        user_id=payment.user_id, 
        amount=payment.amount, 
        description=payment.description
    )

@app.delete("/payments/{payment_id}", response_model=dict)
async def delete_payment(payment_id: int, db: AsyncSession = Depends(get_db)):
    result = await db.execute(select(Payment).where(Payment.id == payment_id))
    payment = result.scalars().first()
    if not payment:
        raise HTTPException(status_code=404, detail="Payment not found")
    await db.delete(payment)
    await db.commit()
    return {"message": "Payment deleted successfully"}