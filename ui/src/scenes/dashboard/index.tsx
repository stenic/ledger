import { Box, Typography } from "@mui/material";
import Header from "../../components/Header";
import { useAuth } from "react-oidc-context";
import { ResponsiveTimeRange } from "@nivo/calendar";
import { useGqlQuery } from "../../utils/http";
import gql from "graphql-tag";
import FlexBetween from "../../components/FlexBetween";
import { nivoTheme } from "../../theme";

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
      // from="2018-04-01"
      // to="2018-08-12"
      emptyColor="#333"
      colors={["#61cdbb", "#97e3d5", "#e8c1a0", "#f47560"]}
      // margin={{ top: 40, right: 40, bottom: 100, left: 40 }}
      dayBorderWidth={0}
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
      </Box>
    </Box>
  );
};

export default Dashboard;
