from quickgrpc.utils import get_stub
import grpc_assets.employeemgmt.employeemgmt_pb2 as pb2
from unittest import TestCase


class EmployeemgmtClientTestCase(TestCase):
    stub = get_stub("employeemgmt")

    def test_crd(self):
        self.stub.create(pb2.EmployeeInp(
                name="ashu",
                dept_id=1,
                designation="dev"
            ))
        print("created employee ashu")
        self.stub.create(pb2.EmployeeInp(
            name="jane",
            dept_id=2,
            designation="devops"
        ))
        print("created employee jane")
        res = self.stub.list(pb2.Empty())
        print(f"listing employees: {res}")
        ashu = self.stub.retrieve(pb2.EmployeeID(id=1))
        print(f"retrieved ashu: {ashu}")
        self.stub.delete(pb2.EmployeeID(id=1))
        print("deleted ashu")
        res = self.stub.list(pb2.Empty())
        print(f"listing employees: {res}")
