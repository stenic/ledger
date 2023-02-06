import { ResponsiveTimeRange } from "@nivo/calendar";
import { nivoTheme } from "../../theme";

export const Timeline = ({
  data,
  clickHandler,
}: {
  data: any;
  clickHandler?: (datum: any, event: React.MouseEvent) => void;
}) => {
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
      onClick={clickHandler}
    />
  );
};
