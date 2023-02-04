import { Box, Typography } from "@mui/material";
import Header from "../../components/Header";
import { useAuth } from "react-oidc-context";
import { ResponsiveTimeRange } from "@nivo/calendar";
import { useGqlQuery } from "../../utils/http";
import gql from "graphql-tag";
import FlexBetween from "../../components/FlexBetween";
import { nivoTheme } from "../../theme";
import { ResponsivePie } from "@nivo/pie";
import { VersionData } from "../../types/version";

import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { useTranslation } from "react-i18next";

const Timeline = ({ data }: { data: any }) => {
  return (
    <ResponsiveTimeRange
      data={data}
      to={new Date().toISOString().substring(0, 10)}
      emptyColor="#333"
      colors={["#61cdbb", "#97e3d5", "#e8c1a0", "#f47560"]}
      // margin={{ top: 40, right: 40, bottom: 100, left: 40 }}
      dayBorderWidth={3}
      dayBorderColor="#1F2A40"
      theme={{
        ...nivoTheme,
        labels: { text: { fill: "#fff" } },
      }}
    />
  );
};

const StatBox = ({
  title,
  value,
  icon,
  description,
}: {
  title: string;
  value: string;
  icon?: any;
  description: string;
}) => {
  return (
    <Box
      gridColumn="span 2"
      gridRow="span 1"
      display="flex"
      flexDirection="column"
      justifyContent="space-between"
      p="1.25 rem 1rem"
      flex="1 1 100%"
    >
      <FlexBetween>
        <Typography variant="h6">{title}</Typography>
        {icon}
      </FlexBetween>
      <Typography variant="h3" fontWeight="600">
        {value}
      </Typography>
      <Typography>{description}</Typography>
    </Box>
  );
};

const Pie = ({ data }: { data: any }) => {
  return (
    <ResponsivePie
      data={data}
      margin={{ top: 40, right: 80, bottom: 80, left: 80 }}
      innerRadius={0.5}
      padAngle={0.7}
      cornerRadius={3}
      activeOuterRadiusOffset={8}
      borderWidth={1}
      arcLinkLabelsSkipAngle={10}
      arcLinkLabelsThickness={2}
      arcLabelsSkipAngle={10}
      theme={{
        ...nivoTheme,
        labels: { text: { fill: "#fff" } },
      }}
    />
  );
};

const LastTable = ({ data }: { data: Array<VersionData> }) => {
  const { t } = useTranslation();

  return (
    <TableContainer sx={{ maxHeight: "100%" }}>
      <Table sx={{ minWidth: 650 }} stickyHeader>
        <TableHead>
          <TableRow>
            <TableCell>{t("table_version_timestamp")}</TableCell>
            <TableCell>{t("table_version_application")}</TableCell>
            <TableCell>{t("table_version_environment")}</TableCell>
            <TableCell>{t("table_version_location")}</TableCell>
            <TableCell align="right">{t("table_version_version")}</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data.map((row: VersionData) => (
            <TableRow
              key={row.id}
              sx={{ "&:last-child td, &:last-child th": { border: 0 } }}
            >
              <TableCell>{row.timestamp}</TableCell>
              <TableCell component="th" scope="row">
                {row.application.name}
              </TableCell>
              <TableCell>{row.environment.name}</TableCell>
              <TableCell>{row.location.name}</TableCell>
              <TableCell align="right">{row.version}</TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

const Dashboard = () => {
  const auth = useAuth();
  const { t } = useTranslation();

  const { data, isLoading } = useGqlQuery(
    ["stats", "version"],
    gql`
      query {
        environments {
          name
        }
        applications {
          name
        }
        locations {
          name
        }
        totalVersions
        versionCountPerDay {
          day: timstamp
          value: count
        }
        lastVersions(days: 7) {
          id
          application {
            name
          }
          location {
            name
          }
          environment {
            name
          }
          version
          timestamp
        }
      }
    `,
    {},
    (d) => {
      let grouped: any = {};
      d.lastVersions.forEach((e: VersionData) => {
        const key: string = [e.location.name, e.environment.name]
          .filter((e) => e.length > 0)
          .join(" / ");
        grouped[key] = grouped[key] || {
          value: 0,
          id: key,
        };
        grouped[key].value++;
      });

      return {
        ...d,
        grouped: Object.values(grouped),
      };
    }
  );

  return (
    <Box m="20px">
      <Header
        title={t("dashboard_title")}
        subtitle={t("dashboard_welcome", {
          username: auth.user?.profile.preferred_username,
        })}
      />

      {isLoading ? (
        <div>{t("app_loading")}</div>
      ) : (
        <Box
          sx={{
            display: "grid",
            gap: "10px",
            gridTemplateColumns: "repeat(12, 1fr)",
            gridAutoRows: "160px",
            "& > div": {
              bgcolor: "background.paper",
              p: 1,
            },
          }}
        >
          <StatBox
            title={t("dashboard_versions_title")}
            description={t("dashboard_versions_description")}
            value={data?.totalVersions}
          />
          <StatBox
            title={t("dashboard_applications_title")}
            description={t("dashboard_applications_description")}
            value={data?.applications.length || "..."}
          />
          <Box
            sx={{
              gridRow: "span 2",
              gridColumn: "span 8",
            }}
          >
            <Timeline data={data.versionCountPerDay} />
          </Box>
          <StatBox
            title={t("dashboard_environments_title")}
            description={t("dashboard_environments_description")}
            value={data?.environments.length || "..."}
          />
          <StatBox
            title={t("dashboard_locations_title")}
            description={t("dashboard_locations_description")}
            value={data?.locations.length || "..."}
          />
          <Box
            sx={{
              gridRow: "span 3",
              gridColumn: "span 7",
              overflow: "hidden",
            }}
          >
            <LastTable data={data.lastVersions} />
          </Box>
          <Box
            sx={{
              gridRow: "span 3",
              gridColumn: "span 5",
            }}
          >
            <Pie data={data.grouped} />
          </Box>
        </Box>
      )}
    </Box>
  );
};

export default Dashboard;
