import { useAuth } from "react-oidc-context";
import { GraphQLClient } from "graphql-request";
import { useQuery } from "@tanstack/react-query";
import { DocumentNode } from "graphql/language/ast";

export const useGqlFetch = () => {
  const auth = useAuth();
  const token = auth.user?.access_token;
  return <T>(query: string): Promise<T> => {
    return fetch("/query", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({ query }),
    }).then((response) => {
      if (response.status === 401) {
        auth.signinSilent();
      }
      if (!response.ok) {
        throw new Error(response.statusText);
      }
      return response.json();
    });
  };
};

const rawResponseCb = (data: any) => data;

export const useGqlQuery = (
  key: string[],
  query: DocumentNode,
  variables: object = {},
  resultCallback: (dataa: any) => any = rawResponseCb,
  config: object = {}
) => {
  const token = useAuth().user?.access_token;
  const headers = {
    headers: {
      authorization: `Bearer ${token}`,
    },
  };

  const graphQLClient = new GraphQLClient("/query", headers);
  const fetchData = async () =>
    await graphQLClient.request(query, variables).then(resultCallback);

  return useQuery({
    queryKey: key,
    queryFn: fetchData,
    refetchOnWindowFocus: false,
    ...config,
  });
};
