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
        routes: {
          "GET /": "./api/health/main.go",
          "POST /auth/create_account": "./api/auth/create_account/main.go",
        },
        customDomain: isProd ? domain : undefined,
      });
      return api;
}