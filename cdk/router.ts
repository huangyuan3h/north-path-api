import { Certificate } from "aws-cdk-lib/aws-certificatemanager";
import { Api, Stack } from "sst/constructs";
const certArn =
  "arn:aws:acm:ap-southeast-1:319653899185:certificate/85267517-ebb2-42d2-83f9-22efbadd5c8f";

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
            allowOrigins: ["https://www.north-path.site","http://localhost:3000"],
            allowCredentials: true,
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