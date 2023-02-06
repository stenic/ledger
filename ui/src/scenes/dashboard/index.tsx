import { Box } from "@mui/material";
import Header from "../../components/Header";
import { useAuth } from "react-oidc-context";
import { useGqlQuery } from "../../utils/http";
import gql from "graphql-tag";
import { VersionData } from "../../types/version";
import { useTranslation } from "react-i18next";
import { LocationData } from "../../types/location";
import { EnvironmentData } from "../../types/environment";
import { ApplicationData } from "../../types/application";
import { Filter } from "./filters";
import { Pie } from "./pie";
import { Timeline } from "./timeline";
import { StatBox } from "./statbox";
import { LastTable } from "./table";
import { useUrlSearchParams } from "use-url-search-params";

const Dashboard = () => {
  const auth = useAuth();
  const { t } = useTranslation();

  // const [dateFilter, setDateFilter] = useHashParam("dateFilter", "");
  const [params, setParams] = useUrlSearchParams({
    location: "",
    environment: "",
    application: "",
  });

  // const dateFilterHandler = (datum: any, event: any) => {
  //   setDateFilter(datum.value ? datum.day : undefined);
  // };

  const filterHandler = (filters: {
    location: string;
    environment: string;
    application: string;
  }) => {
    setParams({ ...params, ...filters });
  };

  const { data, isLoading } = useGqlQuery(
    [
      "stats",
      "version",
      `${params.location}`,
      `${params.environment}`,
      `${params.application}`,
    ],
    gql`
      query ($appFilter: String, $envFilter: String, $locFilter: String) {
        environments(
          filter: {
            application: $appFilter
            environment: $envFilter
            location: $locFilter
          }
        ) {
          name
        }
        applications(
          filter: {
            application: $appFilter
            environment: $envFilter
            location: $locFilter
          }
        ) {
          name
        }
        locations(
          filter: {
            application: $appFilter
            environment: $envFilter
            location: $locFilter
          }
        ) {
          name
        }
        totalVersions(
          filter: {
            application: $appFilter
            environment: $envFilter
            location: $locFilter
          }
        )
        versionCountPerDay(
          filter: {
            application: $appFilter
            environment: $envFilter
            location: $locFilter
          }
        ) {
          day: timstamp
          value: count
        }
        lastVersions(days: 120) {
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
    {
      envFilter: params.environment,
      appFilter: params.application,
      locFilter: params.location,
    },
    (d) => {
      let grouped: any = {};
      // if (dateFilter) {
      //   d.lastVersions = d.lastVersions.filter((e: VersionData) => {
      //     return e.timestamp.substring(0, 10) === dateFilter;
      //   });
      // }
      if (params.location) {
        d.lastVersions = d.lastVersions.filter((e: VersionData) => {
          return e.location.name.match(`${params.location}`);
        });
      }
      if (params.application) {
        d.lastVersions = d.lastVersions.filter((e: VersionData) => {
          return e.application.name.match(`${params.application}`);
        });
      }
      if (params.environment) {
        d.lastVersions = d.lastVersions.filter((e: VersionData) => {
          return e.environment.name.match(`${params.environment}`);
        });
      }
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
      >
        <Filter
          locations={data?.locations.map((l: LocationData) => l.name)}
          environments={data?.environments.map((l: EnvironmentData) => l.name)}
          applications={data?.applications.map((l: ApplicationData) => l.name)}
          initialState={{
            location: `${params.location}`,
            environment: `${params.environment}`,
            application: `${params.application}`,
          }}
          filterCallback={filterHandler}
        />
      </Header>

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
            <Timeline
              data={data.versionCountPerDay}
              // clickHandler={dateFilterHandler}
            />
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
