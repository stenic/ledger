import { Box } from "@mui/material";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";
import Header from "../../components/Header";
import { useGqlQuery } from "../../utils/http";
import gql from "graphql-tag";
import { VersionData } from "../../types/version";

const LastVersions = () => {
  const { data, isLoading } = useGqlQuery(
    ["lastversions", "version"],
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
          timestamp
          environment {
            name
          }
          version
        }
      }
    `
  );

  const columns = [
    {
      field: "application",
      headerName: "Application",
      flex: 1,
      cellClassName: "name-column--cell",
      valueGetter: (params: { row: VersionData }) =>
        params.row.application.name,
    },
    {
      field: "location",
      headerName: "Location",
      flex: 1,
      valueGetter: (params: { row: VersionData }) => params.row.location.name,
    },
    {
      field: "environment",
      headerName: "Environment",
      flex: 1,
      valueGetter: (params: { row: VersionData }) =>
        params.row.environment.name,
    },
    {
      field: "version",
      headerName: "Version",
      flex: 1,
      // renderCell: (params) => (
      //   <Typography color="#4cceac">{params.row.cost}</Typography>
      // ),
    },
    {
      field: "timestamp",
      headerName: "Date",
      flex: 1,
      type: "dateTime",
      valueGetter: ({ value }: { value: string }) => value && new Date(value),
    },
  ];

  return (
    <Box m="20px">
      <Header title="Versions" subtitle="List of versions" />
      <Box
        m="40px 0 0 0"
        height="75vh"
        sx={{
          "& .MuiDataGrid-root": {
            border: "none",
            color: "#fff",
          },
          "& .MuiDataGrid-cell": {
            borderBottom: "none",
          },
          "& .name-column--cell": {
            color: "#94e2cd",
          },
          "& .MuiDataGrid-columnHeaders": {
            backgroundColor: "#3e4396",
            borderBottom: "none",
          },
          "& .MuiDataGrid-virtualScroller": {
            backgroundColor: "#1F2A40",
          },
          "& .MuiDataGrid-footerContainer": {
            borderTop: "none",
            backgroundColor: "#3e4396",
          },
          "& .MuiCheckbox-root": {
            color: `#b7ebde !important`,
          },
          "& .MuiButtonBase-root": {
            color: "#fff !important",
          },
        }}
      >
        <DataGrid
          loading={isLoading}
          rows={data?.lastVersions ? data.lastVersions : []}
          columns={columns}
          components={{
            Toolbar: GridToolbar,
          }}
        />
      </Box>
    </Box>
  );
};

export default LastVersions;
