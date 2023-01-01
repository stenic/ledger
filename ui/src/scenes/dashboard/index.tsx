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
import Paper from "@mui/material/Paper";

const Timeline = () => {
  const { data, isLoading } = useGqlQuery(
    ["versiontimeline"],
    gql`
      query {
        versionCountPerDay {
          day: timstamp
          value: count
        }
      }
    `
  );

  if (isLoading) return <div>Loading</div>;

  return (
    <ResponsiveTimeRange
      data={data?.versionCountPerDay}
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

const Pie = () => {
  const { data, isLoading } = useGqlQuery(
    ["piedata"],
    gql`
      query {
        lastVersions {
          application {
            name
          }
          location {
            name
          }
          environment {
            name
          }
        }
      }
    `,
    {},
    (d) => {
      let grouped: any = {};
      const groupBy = "environment";
      d.lastVersions.forEach((e: VersionData) => {
        grouped[e[groupBy].name] = grouped[e[groupBy].name] || {
          value: 0,
          id: e[groupBy].name,
        };
        grouped[e[groupBy].name].value++;
      });

      return {
        ...d,
        grouped: Object.values(grouped),
      };
    }
  );

  if (isLoading) return <div>Loading</div>;

  return (
    <ResponsivePie
      data={data.grouped}
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

const LastTable = () => {
  const { data, isLoading } = useGqlQuery(
    ["tabledata"],
    gql`
      query {
        lastVersions {
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
    `
  );

  if (isLoading) return <div>Loading</div>;

  return (
    <TableContainer>
      <Table sx={{ minWidth: 650 }} aria-label="simple table">
        <TableHead>
          <TableRow>
            <TableCell>Timestamp</TableCell>
            <TableCell>Application</TableCell>
            <TableCell>Environment</TableCell>
            <TableCell>Location</TableCell>
            <TableCell align="right">Version</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {data.lastVersions.map((row: VersionData) => (
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

  const { data } = useGqlQuery(
    ["stats"],
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
      }
    `
  );

  return (
    <Box m="20px">
      <Header
        title="Dashboard"
        subtitle={`Hello ${auth.user?.profile.preferred_username}, welcome to your dashboard`}
      />

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
          title="Total"
          value={data?.totalVersions}
          description="Total number of versions"
        />
        <StatBox
          title="Applications"
          value={data?.applications.length || "..."}
          description="Unique applications"
        />
        <Box
          sx={{
            gridRow: "span 2",
            gridColumn: "span 8",
          }}
        >
          <Timeline />
        </Box>
        <StatBox
          title="Environments"
          value={data?.environments.length || "..."}
          description="Unique environments"
        />
        <StatBox
          title="Locations"
          value={data?.locations.length || "..."}
          description="Unique locations"
        />
        <Box
          sx={{
            gridRow: "span 3",
            gridColumn: "span 7",
            overflow: "scroll",
          }}
        >
          <LastTable />
        </Box>
        <Box
          sx={{
            gridRow: "span 3",
            gridColumn: "span 5",
          }}
        >
          <Pie />
        </Box>
      </Box>
    </Box>
  );
};

export default Dashboard;
