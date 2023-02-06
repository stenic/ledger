import { VersionData } from "../../types/version";
import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";
import { useTranslation } from "react-i18next";

export const LastTable = ({ data }: { data: Array<VersionData> }) => {
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
