from grpc_assets.employeemgmt import employeemgmt_pb2
from grpc_assets.employeemgmt.employeemgmt_pb2_grpc import EmployeemgmtServicer

DB: dict = {}

class EmployeemgmtService(EmployeemgmtServicer):
    def create(self, employee_inp, context):
        employee = employeemgmt_pb2.Employee(
            name=employee_inp.name,
            dept_id=employee_inp.dept_id,
            designation=employee_inp.designation
        )
        DB[len(DB) + 1] = employee
        return employee

    def list(self, _, context):
        res = employeemgmt_pb2.Employees()
        res.employees.extend([employee for _, employee in DB.items()])
        return res

    def retrieve(self, employee_id, context):
        return DB[employee_id.id]

    def delete(self, employee_id, context):
        del DB[employee_id.id]
        return employeemgmt_pb2.Empty()
