import { SSTConfig } from "sst";
import getApi from "./cdk/router";


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

     const api = getApi(stack)
      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;
