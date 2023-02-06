import { nivoTheme } from "../../theme";
import { ResponsivePie } from "@nivo/pie";

export const Pie = ({ data }: { data: any }) => {
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
