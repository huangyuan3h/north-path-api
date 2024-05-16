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
          status: "string",
          postId: "string",
          email:"string",
          category: "string",
          subject: "string",
          location: "string",
          content: "string",
          images: "string",
          topic:"string",
          createdDate: "string",
          updatedDate: "string",
          like:"number",
          sortingScore:"number",
        },
        primaryIndex: { partitionKey: "postId"},
        globalIndexes: { "all": { partitionKey: "status", sortKey: "sortingScore" },
        "myPost": { partitionKey: "email", sortKey: "updatedDate" }, 
        "category": { partitionKey: "category", sortKey: "sortingScore" },
        "location": { partitionKey: "location", sortKey: "sortingScore" }, }
      });



      return {authTable, userTable, postTable}
}
