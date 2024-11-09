import json
import pika
from dotenv import dotenv_values
from fastapi import FastAPI, HTTPException, Request
from pydantic import BaseModel
from typing import List, Optional

app = FastAPI()

class Employee(BaseModel):
    id: Optional[int] = None
    name: str

employees = [
    {"id": 1, "name": "Ashley"},
    {"id": 2, "name": "Kate"},
    {"id": 3, "name": "Joe"}
]

nextEmployeeId = 4

def get_employee(id: int):
    return next((e for e in employees if e['id'] == id), None)

def employee_is_valid(employee: dict):
    return all(key == "name" for key in employee.keys())

@app.get("/ping")
async def ping():
    return "pong 🏓"

@app.get("/employees", response_model=List[Employee])
async def get_employees():
    return employees

@app.get("/employees/{id}", response_model=Employee)
async def get_employee_by_id(id: int):
    employee = get_employee(id)
    if employee is None:
        raise HTTPException(status_code=404, detail="Employee does not exist")
    return employee

@app.post("/employees", status_code=201)
async def create_employee(request: Request):
    global nextEmployeeId
    employee = await request.json()
    if not employee_is_valid(employee):
        raise HTTPException(status_code=400, detail="Invalid employee properties.")
    
    employee['id'] = nextEmployeeId
    nextEmployeeId += 1
    employees.append(employee)
    
    return {"location": f"/employees/{employee['id']}"}

@app.put("/employees/{id}", response_model=Employee)
async def update_employee(id: int, request: Request):
    employee = get_employee(id)
    if employee is None:
        raise HTTPException(status_code=404, detail="Employee does not exist")
    
    updated_employee = await request.json()
    if not employee_is_valid(updated_employee):
        raise HTTPException(status_code=400, detail="Invalid employee properties.")
    
    employee.update(updated_employee)
    return employee

@app.delete("/employees/{id}", response_model=Employee)
async def delete_employee(id: int):
    global employees
    employee = get_employee(id)
    if employee is None:
        raise HTTPException(status_code=404, detail="Employee does not exist")
    
    employees = [e for e in employees if e['id'] != id]
    return employee

if __name__ == "__main__":
    import uvicorn
    config = dotenv_values(".env")
    connection = pika.BlockingConnection(pika.ConnectionParameters(config.get("QUEUE_HOST")))
    channel = connection.channel()
    uvicorn.run(app, host="0.0.0.0", port=8002)
