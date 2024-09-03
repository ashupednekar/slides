from grpc_assets.employeemgmt import employeemgmt_pb2
from grpc_assets.employeemgmt.employeemgmt_pb2_grpc import EmployeemgmtServicer

DB: dict = {}

class EmployeemgmtService(EmployeemgmtServicer):
    def create(self, employee_inp, context):
        employee = {
            "name": employee_inp.name,
            "dept_id": employee_inp.dept_id,
            "designation": employee_inp.designation
        } 
        DB[len(DB) + 1] = employee
        return employeemgmt_pb2.Employee(
            name=employee["name"],
            dept_id=employee["dept_id"],
            designation=employee["designation"]
        )

    def list(self, _, context):
        res = employeemgmt_pb2.Employees()
        res.employees.extend([employeemgmt_pb2.Employee(
            name=entry["name"],
            dept_id=entry["dept_id"],
            designation=entry["designation"]
        ) for entry in DB.values()])
        return res

    def update(self, employee_inp, context):
        employee = {
            "name": employee_inp.name,
            "dept_id": employee_inp.dept_id,
            "designation": employee_inp.designation
        } 
        #DB[] = employee
        return employeemgmt_pb2.Employee(
            name=employee["name"],
            dept_id=employee["dept_id"],
            designation=employee["designation"]
        )

    def retrieve(self, employee_id, context):
        employee = DB[employee_id.id]
        return employeemgmt_pb2.Employee(
            name=employee["name"],
            dept_id=employee["dept_id"],
            designation=employee["designation"]
        )


    def delete(self, employee_id, context):
        del DB[employee_id.id]
        return employeemgmt_pb2.Empty()
