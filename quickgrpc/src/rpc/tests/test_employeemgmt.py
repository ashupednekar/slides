from quickgrpc.utils import get_stub
import grpc_assets.employeemgmt.employeemgmt_pb2 as pb2
from unittest import TestCase

class EmployeemgmtClientTestCase(TestCase):
    stub = get_stub("employeemgmt")

    def test_crud(self):
        # stub.<your rpc method>(pb2.<your proto message>(params))
        r1 = self.stub.create(pb2.EmployeeInp(
            name="ashu",
            dept_id=1,
            designation="dev"
        ))
        print(f"created employee: {r1}")
        r2 = self.stub.create(pb2.EmployeeInp(
            name="jane",
            dept_id=1,
            designation="devops"
        ))
        print(f"created employee: {r2}")
        l = self.stub.list(pb2.Empty())
        print(f"employee list\n: {l}")
        print(self.stub.retrieve(pb2.EmployeeID(id=1)))
        print(self.stub.delete(pb2.EmployeeID(id=1)))
        print(self.stub.list(pb2.Empty()))

   
