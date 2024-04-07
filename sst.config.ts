import { SSTConfig } from "sst";
import getApi from "./cdk/router";
import { getTableConfig } from "./cdk/table";


export default {
  config(_input) {
    return {
      name: "north-path-api",
      region: "us-east-1",
    };
  },
  stacks(app) {
    app.setDefaultFunctionProps({
      runtime: "go",
    });
    
    app.stack(function Stack({ stack }) {

      const {authTable, userTable} = getTableConfig(stack);

      const api = getApi(stack);

      api.attachPermissions([authTable, userTable]);

      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;

