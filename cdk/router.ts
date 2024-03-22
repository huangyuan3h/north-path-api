import { Certificate } from "aws-cdk-lib/aws-certificatemanager";
import { Api, Stack } from "sst/constructs";
const certArn =
  "arn:aws:acm:ap-southeast-1:319653899185:certificate/e09b910f-b3b6-47d3-a50c-473ae799d532";

export default (stack:Stack)=>{
    const isProd = stack.stage === "prod";

          // define domain
          const domain = {
            domainName: "api.north-path.site",
            isExternalDomain: true,
            cdk: {
              certificate: Certificate.fromCertificateArn(stack, "MyCert", certArn),
            },
          };

    const api = new Api(stack, "api", {
        cors: {
            allowMethods: ["ANY"],
            allowHeaders: ["Authorization"],
            allowOrigins: ["https://www.north-path.site"],
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
        },
        customDomain: isProd ? domain : undefined,
      });
      return api;
}