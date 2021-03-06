import React, { useMemo } from "react";
import { Routes, Route, useParams, Navigate } from "react-router-dom";
import { ApolloProvider } from "@apollo/client";
import { makeClient } from "./graphql/adminapi/apollo";
import ScreenLayout from "./ScreenLayout";

import UsersScreen from "./graphql/adminapi/UsersScreen";
import AddUserScreen from "./graphql/adminapi/AddUserScreen";
import UserDetailsScreen from "./graphql/adminapi/UserDetailsScreen";
import AddEmailScreen from "./graphql/adminapi/AddEmailScreen";
import AddPhoneScreen from "./graphql/adminapi/AddPhoneScreen";
import AddUsernameScreen from "./graphql/adminapi/AddUsernameScreen";
import ResetPasswordScreen from "./graphql/adminapi/ResetPasswordScreen";

import AuthenticationConfigurationScreen from "./graphql/portal/AuthenticationConfigurationScreen";
import AnonymousUsersConfigurationScreen from "./graphql/portal/AnonymousUsersConfigurationScreen";
import SingleSignOnConfigurationScreen from "./graphql/portal/SingleSignOnConfigurationScreen";
import PasswordsScreen from "./graphql/portal/PasswordsScreen";
import PasswordlessAuthenticatorScreen from "./graphql/portal/PasswordlessAuthenticatorScreen";
import OAuthClientConfigurationScreen from "./graphql/portal/OAuthClientConfigurationScreen";
import CreateOAuthClientScreen from "./graphql/portal/CreateOAuthClientScreen";
import EditOAuthClientScreen from "./graphql/portal/EditOAuthClientScreen";
import UserInterfaceScreen from "./graphql/portal/UserInterfaceScreen";
import DNSConfigurationScreen from "./graphql/portal/DNSConfigurationScreen";
import SettingsScreen from "./graphql/portal/SettingsScreen";

const AppRoot: React.FC = function AppRoot() {
  const { appID } = useParams();
  const client = useMemo(() => {
    return makeClient(appID);
  }, [appID]);
  return (
    <ApolloProvider client={client}>
      <ScreenLayout>
        <Routes>
          <Route path="/" element={<Navigate to="users/" replace={true} />} />
          <Route path="/users/" element={<UsersScreen />} />
          <Route path="/users/add-user/" element={<AddUserScreen />} />
          <Route
            path="/users/:userID/"
            element={<Navigate to="details/" replace={true} />}
          />
          <Route
            path="/users/:userID/details/"
            element={<UserDetailsScreen />}
          />
          <Route
            path="/users/:userID/details/add-email"
            element={<AddEmailScreen />}
          />
          <Route
            path="/users/:userID/details/add-phone"
            element={<AddPhoneScreen />}
          />
          <Route
            path="/users/:userID/details/add-username"
            element={<AddUsernameScreen />}
          />
          <Route
            path="/users/:userID/details/reset-password"
            element={<ResetPasswordScreen />}
          />
          <Route
            path="/configuration/authentication/"
            element={<AuthenticationConfigurationScreen />}
          />
          <Route
            path="/configuration/anonymous-users"
            element={<AnonymousUsersConfigurationScreen />}
          />
          <Route
            path="/configuration/single-sign-on"
            element={<SingleSignOnConfigurationScreen />}
          />
          <Route
            path="/configuration/passwords"
            element={<PasswordsScreen />}
          />
          <Route
            path="/configuration/passwordless-authenticator"
            element={<PasswordlessAuthenticatorScreen />}
          />
          <Route
            path="/configuration/oauth-clients"
            element={<OAuthClientConfigurationScreen />}
          />
          <Route
            path="/configuration/oauth-clients/add"
            element={<CreateOAuthClientScreen />}
          />
          <Route
            path="/configuration/oauth-clients/:clientID/edit"
            element={<EditOAuthClientScreen />}
          />
          <Route
            path="/configuration/user-interface"
            element={<UserInterfaceScreen />}
          />
          <Route
            path="/configuration/dns"
            element={<DNSConfigurationScreen />}
          />
          <Route path="/configuration/settings" element={<SettingsScreen />} />
        </Routes>
      </ScreenLayout>
    </ApolloProvider>
  );
};

export default AppRoot;
