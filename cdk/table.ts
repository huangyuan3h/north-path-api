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

      const userTable = new Table(stack, "user", {
        fields: {
          email: "string",
          avatar: "string",
          userName: "string",
          bio: "string"
        },
        primaryIndex: { partitionKey: "email" },
      });

      const postTable = new Table(stack, "posts", {
        fields: {
          postId: "string",
          email:"string",
          subject: "string",
          category: "string",
          content: "string",
          createdDate: "number",
        },
        primaryIndex: { partitionKey: "postId" },
      });



      return {authTable, userTable, postTable}
}
