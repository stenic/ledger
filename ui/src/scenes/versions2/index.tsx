import { Box } from "@mui/material";
import Header from "../../components/Header";
import ReactDataGrid from "@inovua/reactdatagrid-community";
// import "@inovua/reactdatagrid-community/index.css";
import "@inovua/reactdatagrid-community/base.css";
import "@inovua/reactdatagrid-community/theme/default-dark.css";

import { useGqlFetch } from "../../utils/http";
import { useCallback } from "react";
import { VersionData } from "../../types/version";
import {
  TypeFilterValue,
  TypeSortInfo,
} from "@inovua/reactdatagrid-community/types";

const gridStyle = { minHeight: 600, marginTop: 10 };

// const defaultFilterValue = [
//   { name: "environment", type: "string", operator: "contains", value: "" },
//   { name: "application", type: "string", operator: "contains", value: "" },
//   { name: "version", type: "string", operator: "contains", value: "" },
// ];

const columns = [
  {
    name: "id",
    header: "Id",
    defaultVisible: false,
    type: "number",
    defaultWidth: 60,
  },
  { name: "environment", header: "Environment", defaultFlex: 1 },
  { name: "application", header: "Application", defaultFlex: 1 },
  { name: "version", header: "Version", groupBy: false, defaultFlex: 1 },
  { name: "timestamp", header: "Timestamp", groupBy: false, defaultFlex: 1 },
];

const Versions2 = () => {
  const gqlFetch = useGqlFetch();

  const handleSortChange = (sortInfo: TypeSortInfo) => {};
  const handleFilterValueChanged = (filterValue: TypeFilterValue) => {};

  const loadData = async () => {
    return gqlFetch<{ data: { versions: VersionData[] } }>(
      "{ versions { id, application { name }, timestamp, environment { name }, version } }"
    ).then((d) => {
      const result: any[] = d.data.versions.map((v) => ({
        ...v,
        application: v.application.name,
        environment: v.environment.name,
      }));
      return { data: result, count: result.length };
    });
  };

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const dataSource = useCallback(loadData, []);

  return (
    <Box m="20px">
      <Header title="Versions2" subtitle="List of versions" />
      <ReactDataGrid
        idProperty="id"
        style={gridStyle}
        columns={columns}
        // defaultFilterValue={defaultFilterValue}
        defaultGroupBy={[]}
        dataSource={dataSource}
        onSortInfoChange={handleSortChange}
        onFilterValueChange={handleFilterValueChanged}
      />
    </Box>
  );
};

export default Versions2;
