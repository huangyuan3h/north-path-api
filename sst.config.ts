import { SSTConfig } from "sst";
import { Api } from "sst/constructs";
import { Certificate } from "aws-cdk-lib/aws-certificatemanager";

const certArn =
  "arn:aws:acm:ap-southeast-1:319653899185:certificate/e09b910f-b3b6-47d3-a50c-473ae799d532";

export default {
  config(_input) {
    return {
      name: "north-path-api",
      region: "ap-southeast-1",
    };
  },
  stacks(app) {
    app.setDefaultFunctionProps({
      runtime: "go",
    });
    app.stack(function Stack({ stack }) {
      // define domain
      const domain = {
        domainName: "api.north-path.site",
        isExternalDomain: true,
        cdk: {
          certificate: Certificate.fromCertificateArn(stack, "MyCert", certArn),
        },
      };

      const isProd = stack.stage === "prod";

      const api = new Api(stack, "api", {
        routes: {
          "GET /": "./lambda/main.go",
        },
        customDomain: isProd ? domain : undefined,
      });
      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;
