import { useMemo } from "react";
import { gql, QueryResult, useQuery } from "@apollo/client";

import { client } from "../apollo";
import { PortalAPIEmailAndSmsTemplates } from "../../../types";
import { AppConfigQuery_node_App } from "./__generated__/AppConfigQuery";
import {
  AppAndEmailSmsTemplatesQuery,
  AppAndEmailSmsTemplatesQueryVariables,
} from "./__generated__/AppAndEmailSmsTemplatesQuery";

export const andAndEmailSmsTemplatesQuery = gql`
  query AppAndEmailSmsTemplatesQuery(
    $id: ID!
    $emailHtmlPath: String!
    $emailTextPath: String!
    $smsTextPath: String!
  ) {
    node(id: $id) {
      __typename
      ... on App {
        id
        rawAppConfig
        effectiveAppConfig
        emailHtml: rawConfigFile(path: $emailHtmlPath)
        emailText: rawConfigFile(path: $emailTextPath)
        smsText: rawConfigFile(path: $smsTextPath)
      }
    }
  }
`;

export interface AppAndEmailSmsTemplatesQueryResult
  extends Pick<
    QueryResult<
      AppAndEmailSmsTemplatesQuery,
      AppAndEmailSmsTemplatesQueryVariables
    >,
    "loading" | "error" | "refetch"
  > {
  rawAppConfig: AppConfigQuery_node_App["rawAppConfig"] | null;
  effectiveAppConfig: AppConfigQuery_node_App["effectiveAppConfig"] | null;
  emailAndSmsTemplates: PortalAPIEmailAndSmsTemplates | null;
}

export const useAppAndEmailSmsTemplatesQuery = (
  appID: string,
  emailHtmlTemplatePath: string,
  emailTextTemplatePath: string,
  smsTextTemplatePath: string
): AppAndEmailSmsTemplatesQueryResult => {
  const { data, loading, error, refetch } = useQuery<
    AppAndEmailSmsTemplatesQuery,
    AppAndEmailSmsTemplatesQueryVariables
  >(andAndEmailSmsTemplatesQuery, {
    client,
    variables: {
      id: appID,
      emailHtmlPath: emailHtmlTemplatePath,
      emailTextPath: emailTextTemplatePath,
      smsTextPath: smsTextTemplatePath,
    },
  });

  const queryData = useMemo(() => {
    const appConfigNode = data?.node?.__typename === "App" ? data.node : null;
    return {
      rawAppConfig: appConfigNode?.rawAppConfig ?? null,
      effectiveAppConfig: appConfigNode?.effectiveAppConfig ?? null,
      emailAndSmsTemplates: {
        emailHtml: appConfigNode?.emailHtml,
        emailText: appConfigNode?.emailText,
        smsText: appConfigNode?.smsText,
      },
    };
  }, [data]);

  return { ...queryData, loading, error, refetch };
};
