import { SSTConfig } from "sst";
import getApi from "./cdk/router";
import { getTableConfig } from "./cdk/table";


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

      const {authTable} = getTableConfig(stack)

      const api = getApi(stack)

      api.attachPermissions([authTable]);

      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;

