import { Api, Stack } from "sst/constructs";


export default (stack:Stack)=>{
   
    const api = new Api(stack, "api", {
        cors: {
            allowOrigins: ["https://www.north-path.site","http://localhost:3000", "https://north-path.it-t.xyz"],
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
          "POST /auth/login/google": {
            function: {
              handler:"./api/auth/google/main.go",
              timeout: 10,
              environment: { 
                AVATAR_BUCKET_NAME: process.env.AVATAR_BUCKET_NAME ?? "",
                GOOGLE_CLIENT_ID: process.env.GOOGLE_CLIENT_ID ?? "",
                JWT_SECRET: process.env.JWT_SECRET ?? ""
              },
            }
          },
          // // rcic part deprecated
          // "POST /rcic/search": {
          //   function: {
          //     handler:"./api/rcic/search/main.go",
          //     timeout: 10,
          //   }
          // },
          // posts related part
          "POST /post/createOrUpdate": {
            function: {
              handler:"./api/post/create/main.go",
              timeout: 10,
              environment: { 
                JWT_SECRET: process.env.JWT_SECRET ?? ""
              },
            }
          },
          "GET /post/{id}": {
            function: {
              handler:"./api/post/view/main.go",
              timeout: 10,
            }
          },
          "POST /posts": {
            function: {
              handler:"./api/post/search/main.go",
              timeout: 10,
            }
          },
          // my profile
          "GET /my/profile": {
            function: {
              handler:"./api/user/profile/main.go",
              timeout: 10,
              environment: { 
                JWT_SECRET: process.env.JWT_SECRET ?? ""
              },
            }
          },
          "POST /my/profile": {
            function: {
              handler:"./api/user/update-profile/main.go",
              timeout: 10,
              environment: { 
                JWT_SECRET: process.env.JWT_SECRET ?? ""
              },
            }
          },
          // my posts
          "POST /my/posts": {
            function: {
              handler:"./api/post/my_posts/main.go",
              timeout: 10,
              environment: { 
                JWT_SECRET: process.env.JWT_SECRET ?? ""
              },
            }
          },
          // delete my post
          "POST /my/post/delete": {
            function: {
              handler:"./api/post/delete/main.go",
              timeout: 10,
              environment: { 
                JWT_SECRET: process.env.JWT_SECRET ?? "",
                POST_IMAGE_BUCKET_NAME: process.env.POST_IMAGE_BUCKET_NAME ?? ""
              },
            }
          },
          // send message
          "POST /message/send": {
            function: {
              handler:"./api/message/sendMessage/main.go",
              timeout: 10,
              environment: { 
                EmailToken: process.env.EmailToken ?? ""
              },
            }
          },
        },
      });
      return api;
}