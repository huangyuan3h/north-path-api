import { Api, Stack } from "sst/constructs";


export default (stack:Stack)=>{
   
    const api = new Api(stack, "api", {
        cors: {
            allowOrigins: ["https://www.north-path.site","http://localhost:3000"],
            allowCredentials: true,
            allowHeaders: ["Authorization"],
          },
        routes: {
          "GET /": "./api/health/main.go",
          "POST /auth/create_account": {
            function: {
              handler:"./api/auth/create_account/main.go",
              timeout: 10,
              environment: { AUTH_SECRET: process.env.AUTH_SECRET ?? ""  },
            }
          },
          "POST /auth/login": {
            function: {
              handler:"./api/auth/login/main.go",
              timeout: 10,
              environment: { 
                AUTH_SECRET: process.env.AUTH_SECRET ?? "",
                JWT_SECRET: process.env.JWT_SECRET ?? ""
              },
            }
          },
          // rcic part
          "POST /rcic/search": {
            function: {
              handler:"./api/rcic/search/main.go",
              timeout: 10,
            }
          },
          // posts related part
          "POST /post/create": {
            function: {
              handler:"./api/post/create/main.go",
              timeout: 10,
            }
          },
          "GET /post/{id}": {
            function: {
              handler:"./api/post/view/main.go",
              timeout: 10,
            }
          },
        },
      });
      return api;
}