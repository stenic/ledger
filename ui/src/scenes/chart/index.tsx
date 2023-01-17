import { Box } from "@mui/material";
import Header from "../../components/Header";
import React, { useState } from "react";
import { ResponsiveBar } from "@nivo/bar";
import { nivoTheme } from "../../theme";
import InputLabel from "@mui/material/InputLabel";
import MenuItem from "@mui/material/MenuItem";
import FormControl from "@mui/material/FormControl";
import Select, { SelectChangeEvent } from "@mui/material/Select";
import { useGqlQuery } from "../../utils/http";
import gql from "graphql-tag";
import { VersionData } from "../../types/version";

const Chart = () => {
  const transformCb = (d: { versions: VersionData[] }) => {
    const getDates = (startDate: Date, stopDate: Date) => {
      var dates: any = {};
      var currentDate = startDate;
      while (currentDate <= stopDate) {
        const key = currentDate
          .toLocaleString("sv", { timeZoneName: "short" })
          .substring(0, 10);
        dates[key] = { key };
        currentDate = addDays(currentDate, 1);
      }
      return dates;
    };

    const minDate = d.versions.reduce(
      (prev: VersionData, curr: VersionData) => {
        return prev.timestamp < curr.timestamp ? prev : curr;
      }
    ).timestamp;
    const maxDate = d.versions.reduce(
      (prev: VersionData, curr: VersionData) => {
        return prev.timestamp > curr.timestamp ? prev : curr;
      }
    ).timestamp;

    // compute list of unique items
    const envs = [...new Set(d.versions.map((i) => i.environment.name))].filter(
      (n) => n.length > 0
    );
    const apps = [...new Set(d.versions.map((i) => i.application.name))].filter(
      (n) => n.length > 0
    );
    const locs = [...new Set(d.versions.map((i) => i.location.name))].filter(
      (n) => n.length > 0
    );

    // Create empty data list
    let format = getDates(new Date(minDate), new Date(maxDate));

    // group per date and collect
    d.versions.forEach((item) => {
      const dt = new Date(item.timestamp);
      const key = dt
        .toLocaleString("sv", { timeZoneName: "short" })
        .substring(0, 10);

      if (!(key in format)) {
        format[key] = { key };
      }
      item.environment &&
        (format[key][item.environment.name] =
          (format[key][item.environment.name]
            ? format[key][item.environment.name]
            : 0) + 1);

      item.application &&
        (format[key][item.application.name] =
          (format[key][item.application.name]
            ? format[key][item.application.name]
            : 0) + 1);
    });

    // Ensure dates are in order
    format = Object.keys(format)
      .sort()
      .reduce((obj: any, key: string) => {
        obj[key] = format[key];
        return obj;
      }, {});

    return {
      versions: Object.values(format),
      keys: {
        environment: envs,
        application: apps,
        location: locs,
      },
    };
  };

  const { data } = useGqlQuery(
    ["chart", "version"],
    gql`
      query {
        versions {
          id
          application {
            name
          }
          timestamp
          environment {
            name
          }
          location {
            name
          }
          version
        }
      }
    `,
    {},
    transformCb
  );

  const [chartContent, setChartContent] = useState("environment");

  const handleChange = (event: SelectChangeEvent) => {
    console.log(event);
    setChartContent(event.target.value);
  };

  const addDays = (date: Date, days: number) => {
    var newDate = new Date(date.valueOf());
    newDate.setDate(newDate.getDate() + days);
    return newDate;
  };

  return (
    <Box m="20px">
      <Header title="Deploy" subtitle={chartContent + " deploys per day"} />
      <Box height="80vh">
        <Box maxWidth={150}>
          <FormControl fullWidth>
            <InputLabel id="demo-simple-select-label">Group by</InputLabel>
            <Select
              labelId="demo-simple-select-label"
              id="demo-simple-select"
              value={chartContent}
              label="Age"
              onChange={handleChange}
            >
              <MenuItem value="environment">Environment</MenuItem>
              <MenuItem value="application">Application</MenuItem>
              <MenuItem value="location">Location</MenuItem>
            </Select>
          </FormControl>
        </Box>
        {data?.versions && (
          <ResponsiveBar
            theme={nivoTheme}
            colors={{ scheme: "purple_orange" }}
            data={data?.versions}
            keys={data.keys[chartContent]}
            indexBy="key"
            margin={{ top: 50, right: 130, bottom: 150, left: 60 }}
            axisBottom={{ tickRotation: 45 }}
          />
        )}
      </Box>
    </Box>
  );
};

export default Chart;
