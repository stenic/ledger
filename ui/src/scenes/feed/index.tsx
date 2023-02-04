import { Box, Typography } from "@mui/material";
import Header from "../../components/Header";
import { useGqlQuery } from "../../utils/http";
import gql from "graphql-tag";
import * as React from "react";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import Divider from "@mui/material/Divider";
import ListItemText from "@mui/material/ListItemText";
import ListItemAvatar from "@mui/material/ListItemAvatar";
import Avatar from "@mui/material/Avatar";
import LinearProgress from "@mui/material/LinearProgress";
import { VersionData } from "../../types/version";
import { useTranslation, Trans } from "react-i18next";

const Feed = () => {
  const { t } = useTranslation();
  const { data, isLoading } = useGqlQuery(
    ["feed", "version"],
    gql`
      query {
        versions(orderBy: { timestamp: desc }) {
          id
          application {
            name
          }
          location {
            name
          }
          timestamp
          environment {
            name
          }
          version
        }
      }
    `
  );

  return (
    <Box m="20px">
      <Header title={t("feed_title")} subtitle={t("feed_subtitle")} />
      <List sx={{ width: "100%", bgcolor: "background.paper" }}>
        {isLoading && <LinearProgress />}
        {data &&
          data.versions.map((version: VersionData) => (
            <React.Fragment key={version.id}>
              <ListItem alignItems="flex-start">
                <ListItemAvatar>
                  <Avatar
                    alt={
                      version.application.name
                        ? version.application.name
                        : version.environment.name
                    }
                  />
                </ListItemAvatar>

                <ListItemText
                  primary={
                    <React.Fragment>
                      <Typography
                        sx={{ display: "inline" }}
                        component="span"
                        color="text.primary"
                      >
                        {new Date(version.timestamp).toLocaleDateString()}
                      </Typography>{" "}
                      <Typography
                        sx={{ display: "inline" }}
                        component="span"
                        color="text.secondary"
                      >
                        {new Date(version.timestamp).toLocaleTimeString()}
                      </Typography>
                    </React.Fragment>
                  }
                  secondary={
                    <React.Fragment>
                      <Typography
                        sx={{ display: "inline" }}
                        component="span"
                        variant="body2"
                        color="text.primary"
                      >
                        <Trans
                          i18nKey="feed_item_deployed"
                          t={t}
                          values={{
                            item: version.application.name
                              ? version.application.name + ":" + version.version
                              : version.version,
                            target:
                              version.location.name === ""
                                ? version.environment.name
                                : version.location.name +
                                  "/" +
                                  version.environment.name,
                          }}
                          components={{
                            b: (
                              <Typography
                                component="span"
                                sx={{ fontWeight: 600 }}
                              ></Typography>
                            ),
                          }}
                        />
                      </Typography>
                    </React.Fragment>
                  }
                />
              </ListItem>
              <Divider component="li" />
            </React.Fragment>
          ))}
      </List>
    </Box>
  );
};

export default Feed;
