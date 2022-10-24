import { Construct } from "constructs";
import { App, TerraformStack, CloudBackend, NamedCloudWorkspace, TerraformOutput } from "cdktf";
import * as google from '@cdktf/provider-google';

const project_id = 'jubilant-doodle';
const default_location = 'asia-northeast1';

class MyStack extends TerraformStack {
  constructor(scope: Construct, name: string) {
    super(scope, name);

    new google.GoogleProvider(this, 'Google', {
      project: project_id,
    });

    // https://cloud.google.com/run/docs/authenticating/service-to-service?hl=ja#terraform

    const inter_service_account = new google.ServiceAccount(this, 'inter_service_account', {
      accountId: 'interserviceaccount',
    });

    const private_policy = new google.DataGoogleIamPolicy(this, 'private_policy', {
      binding: [{
        role: 'roles/run.invoker',
        members: [`serviceAccount:${inter_service_account.email}`],
      }],
    });

    const back_run = new google.CloudRunService(this, 'back_run', {
      location: default_location,
      name: 'back',
      template: {
        spec: {
          containers: [{
            image: 'asia-northeast1-docker.pkg.dev/jubilant-doodle/jubilant-doodle/back:latest'
          }],
        },
      },
    });

    new google.CloudRunServiceIamPolicy(this, 'back_private', {
      location: default_location,
      project: project_id,
      service: back_run.name,
      policyData: private_policy.policyData,
    });

    const public_policy = new google.DataGoogleIamPolicy(this, 'public_policy', {
      binding: [{
        role: 'roles/run.invoker',
        members: ['allUsers'],
      }],
    });

    const front_allowed_run = new google.CloudRunService(this, 'front_allowed_run', {
      location: default_location,
      name: 'frontallowed',
      template: {
        spec: {
          containers: [{
            env: [{
              name: 'DST_URL',
              value: back_run.status.get(0).url,
            }],
            image: 'asia-northeast1-docker.pkg.dev/jubilant-doodle/jubilant-doodle/front:latest'
          }],
          serviceAccountName: inter_service_account.email,
        },
      },
    });

    new google.CloudRunServiceIamPolicy(this, 'front_allowed_public', {
      location: default_location,
      project: project_id,
      service: front_allowed_run.name,
      policyData: public_policy.policyData,
    });

    const front_unallowed_run = new google.CloudRunService(this, 'front_unallowed_run', {
      location: default_location,
      name: 'frontunallowed',
      template: {
        spec: {
          containers: [{
            env: [{
              name: 'DST_URL',
              value: back_run.status.get(0).url,
            }],
            image: 'asia-northeast1-docker.pkg.dev/jubilant-doodle/jubilant-doodle/front:latest'
          }],
        },
      },
    });

    new google.CloudRunServiceIamPolicy(this, 'front_unallowed_public', {
      location: default_location,
      project: project_id,
      service: front_unallowed_run.name,
      policyData: public_policy.policyData,
    });

    new TerraformOutput(this, 'front_allowed_url', {
      value: front_allowed_run.status.get(0).url,
    });

    new TerraformOutput(this, 'front_unallowed_url', {
      value: front_unallowed_run.status.get(0).url,
    });

    new TerraformOutput(this, 'back_url', {
      value: back_run.status.get(0).url,
    });

  }
}

const app = new App();
const stack = new MyStack(app, "jubilant-doodle");
new CloudBackend(stack, {
  hostname: "app.terraform.io",
  organization: "hsmtkkdefault",
  workspaces: new NamedCloudWorkspace("jubilant-doodle")
});
app.synth();
