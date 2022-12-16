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
import moment from "moment";

const Feed = () => {
  const { data, isLoading } = useGqlQuery(
    ["versions"],
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
      <Header title="Feed" subtitle="Live feed" />
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
                        {moment(version.timestamp).format("llll")}
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
                        Deployed{" "}
                        <Typography component="span" sx={{ fontWeight: 600 }}>
                          {version.application.name
                            ? version.application.name + ":"
                            : ""}
                          {version.version}
                        </Typography>{" "}
                        to{" "}
                        <Typography component="span" sx={{ fontWeight: 600 }}>
                          {version.location.name === ""
                            ? ""
                            : version.location.name + "/"}
                          {version.environment.name}
                        </Typography>
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
