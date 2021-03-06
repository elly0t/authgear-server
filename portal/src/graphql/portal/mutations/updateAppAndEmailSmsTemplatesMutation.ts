import React from "react";
import { useMutation, gql } from "@apollo/client";
import yaml from "js-yaml";

import { client } from "../apollo";
import { PortalAPIAppConfig } from "../../../types";
import { AppConfigFile } from "../__generated__/globalTypes";
import {
  UpdateAppAndEmailSmsTemplatesConfigMutation,
  UpdateAppAndEmailSmsTemplatesConfigMutationVariables,
  UpdateAppAndEmailSmsTemplatesConfigMutation_updateAppConfig,
} from "./__generated__/UpdateAppAndEmailSmsTemplatesConfigMutation";

// relative to project root
const APP_CONFIG_PATH = "./authgear.yaml";

const updateAppAndEmailSmsTemplatesConfigMutation = gql`
  mutation UpdateAppAndEmailSmsTemplatesConfigMutation(
    $appID: ID!
    $updateFiles: [AppConfigFile!]!
    $emailHtmlTemplatePath: String!
    $emailTextTemplatePath: String!
    $smsTextTemplatePath: String!
  ) {
    updateAppConfig(
      input: { appID: $appID, updateFiles: $updateFiles, deleteFiles: [] }
    ) {
      id
      rawAppConfig
      effectiveAppConfig
      emailHtml: rawConfigFile(path: $emailHtmlTemplatePath)
      emailText: rawConfigFile(path: $emailTextTemplatePath)
      smsText: rawConfigFile(path: $smsTextTemplatePath)
    }
  }
`;

export type AppAndEmailSmsTemplatesConfigUpdater = (
  appConfig: PortalAPIAppConfig,
  templateContents: {
    emailHtml?: string;
    emailText?: string;
    smsText?: string;
  }
) => Promise<UpdateAppAndEmailSmsTemplatesConfigMutation_updateAppConfig | null>;

export function useUpdateAppAndEmailSmsTemplatesConfigMutation(
  appID: string,
  emailHtmlTemplatePath: string,
  emailTextTemplatePath: string,
  smsTextTemplatePath: string
): {
  updateAppAndEmailSmsTemplatesConfig: AppAndEmailSmsTemplatesConfigUpdater;
  loading: boolean;
  error: unknown;
} {
  const [mutationFunction, { error, loading }] = useMutation<
    UpdateAppAndEmailSmsTemplatesConfigMutation,
    UpdateAppAndEmailSmsTemplatesConfigMutationVariables
  >(updateAppAndEmailSmsTemplatesConfigMutation, { client });
  const updateAppAndEmailSmsTemplatesConfig = React.useCallback(
    async (
      appConfig: PortalAPIAppConfig,
      templateContents: {
        emailHtml?: string;
        emailText?: string;
        smsText?: string;
      }
    ) => {
      const appConfigYaml = yaml.safeDump(appConfig);

      const updateFiles: AppConfigFile[] = [
        { path: APP_CONFIG_PATH, content: appConfigYaml },
      ];
      if (templateContents.emailHtml) {
        updateFiles.push({
          path: emailHtmlTemplatePath,
          content: templateContents.emailHtml,
        });
      }
      if (templateContents.emailText) {
        updateFiles.push({
          path: emailTextTemplatePath,
          content: templateContents.emailText,
        });
      }
      if (templateContents.smsText) {
        updateFiles.push({
          path: smsTextTemplatePath,
          content: templateContents.smsText,
        });
      }

      const result = await mutationFunction({
        variables: {
          appID,
          updateFiles,
          emailHtmlTemplatePath,
          emailTextTemplatePath,
          smsTextTemplatePath,
        },
      });
      return result.data?.updateAppConfig ?? null;
    },
    [
      appID,
      emailHtmlTemplatePath,
      emailTextTemplatePath,
      smsTextTemplatePath,
      mutationFunction,
    ]
  );
  return { updateAppAndEmailSmsTemplatesConfig, error, loading };
}
