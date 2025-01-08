from pydantic import BaseModel


class PaymentCreate(BaseModel):
    user_id: int
    amount: float
    description: str | None = None

class PaymentResponse(BaseModel):
    id: int
    user_id: int
    amount: float
    description: str | None