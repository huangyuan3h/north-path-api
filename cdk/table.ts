import { Table, Stack } from "sst/constructs";

export const getTableConfig = (stack: Stack) =>{

    const authTable = new Table(stack, "auth", {
        fields: {
          email: "string",
          password: "string",
          status:"string" // sendEmail, actived, deactivated
        },
        primaryIndex: { partitionKey: "email" },
      });


      return {authTable}
}
