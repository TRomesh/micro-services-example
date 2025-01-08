from sqlalchemy import Column, Float, Integer, String
from sqlalchemy.orm import  declarative_base

Base = declarative_base()

class Payment(Base):
    __tablename__ = "payments"

    id = Column(Integer, primary_key=True, index=True)
    user_id = Column(Integer, index=True, nullable=False)
    amount = Column(Float, nullable=False)
    description = Column(String, nullable=True)