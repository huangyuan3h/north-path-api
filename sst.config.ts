import { SSTConfig } from "sst";
import getApi from "./cdk/router";
import { getTableConfig } from "./cdk/table";
import {  Bucket,  } from 'sst/constructs';
import * as s3 from "aws-cdk-lib/aws-s3";

const bucketArn = "arn:aws:s3:::"+process.env.POST_IMAGE_BUCKET_NAME;

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

      const {authTable, userTable, postTable} = getTableConfig(stack);

      const api = getApi(stack);

      api.attachPermissions([authTable, userTable, postTable]);
      const bucket = new Bucket(stack, 'avatar');
      const bucketImage = new Bucket(stack, "Bucket", {
        cdk: {
          bucket: s3.Bucket.fromBucketArn(stack, "IBucket", bucketArn),
        },
      });
      api.bind([bucket, bucketImage])

      stack.addOutputs({
        ApiEndpoint: api.url,
      });
    });
  },
} satisfies SSTConfig;

