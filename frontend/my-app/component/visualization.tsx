"use client";
import { useState, useEffect } from "react";
import {
  BarChart,
  Bar,
  LineChart,
  Line,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
} from "recharts";

export default function ThroughputVisualization({
  filename,
}: {
  filename: string;
}) {
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    // In a real Next.js app, this would load from /public
    fetch(`/${filename}.json`)
      .then((response) => {
        if (!response.ok) {
          throw new Error("Failed to fetch data");
        }
        return response.json();
      })
      .then((throughputData) => {
        // Process the data for visualization
        const localhostData = prepareComparisonData(
          throughputData["l-l"].drops.throughput,
          throughputData["l-l"].nodrops.throughput
        );

        const piRhoData = prepareComparisonData(
          throughputData["pi-rho"].drops.throughput,
          throughputData["pi-rho"].nodrops.throughput
        );

        // Calculate averages for bar charts
        const localhostAvg = {
          withDrops: calculateAverage(throughputData["l-l"].drops.throughput),
          noDrops: calculateAverage(throughputData["l-l"].nodrops.throughput),
        };

        const piRhoAvg = {
          withDrops: calculateAverage(
            throughputData["pi-rho"].drops.throughput
          ),
          noDrops: calculateAverage(
            throughputData["pi-rho"].nodrops.throughput
          ),
        };

        setData({
          localhost: {
            lineData: localhostData,
            averages: [
              { name: "With Drops", value: localhostAvg.withDrops },
              { name: "No Drops", value: localhostAvg.noDrops },
            ],
          },
          piRho: {
            lineData: piRhoData,
            averages: [
              { name: "With Drops", value: piRhoAvg.withDrops },
              { name: "No Drops", value: piRhoAvg.noDrops },
            ],
          },
        });

        setLoading(false);
      })
      .catch((error) => {
        console.error("Error processing data:", error);
        setError(error.message);
        setLoading(false);
      });
  }, []);

  // Function to prepare data for comparison charts
  const prepareComparisonData = (dropsData, noDropsData) => {
    const maxLength = Math.max(dropsData.length, noDropsData.length);
    const result = [];

    for (let i = 0; i < maxLength; i++) {
      result.push({
        index: i,
        withDrops: i < dropsData.length ? dropsData[i] : null,
        noDrops: i < noDropsData.length ? noDropsData[i] : null,
      });
    }

    return result;
  };

  // Function to calculate average
  const calculateAverage = (data) => {
    return data.reduce((sum, value) => sum + value, 0) / data.length;
  };

  if (loading)
    return (
      <div className="flex justify-center items-center h-64">
        Loading throughput data...
      </div>
    );
  if (error)
    return <div className="text-red-500 p-4">Error loading data: {error}</div>;

  return (
    <div className="p-4">
      <h1 className="text-2xl font-bold mb-8 text-center">
        Network Throughput Comparison (MB/s)
      </h1>

      <div className="grid grid-cols-1 gap-12">
        {/* Localhost to Localhost Section */}
        <div className="bg-white p-4 rounded shadow">
          <h2 className="text-xl font-semibold mb-6 text-center">
            Localhost to Localhost Throughput
          </h2>

          {/* Average Comparison */}
          <div className="mb-8">
            <h3 className="text-lg font-medium mb-3">Average Throughput</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={data.localhost.averages}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis
                  label={{ value: "MB/s", angle: -90, position: "insideLeft" }}
                />
                <Tooltip formatter={(value) => value.toFixed(2) + " MB/s"} />
                <Legend />
                <Bar dataKey="value" fill="#8884d8" name="Throughput (MB/s)" />
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Line Chart Comparison */}
          <div>
            <h3 className="text-lg font-medium mb-3">Throughput per Sample</h3>
            <ResponsiveContainer width="100%" height={400}>
              <LineChart data={data.localhost.lineData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="index"
                  label={{
                    value: "Sample #",
                    position: "insideBottom",
                    offset: -5,
                  }}
                />
                <YAxis
                  label={{ value: "MB/s", angle: -90, position: "insideLeft" }}
                />
                <Tooltip
                  formatter={(value) =>
                    value !== null ? value.toFixed(2) + " MB/s" : "N/A"
                  }
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="withDrops"
                  stroke="#ff7300"
                  name="With Drops"
                  dot={{ r: 3 }}
                  strokeWidth={2}
                  connectNulls
                />
                <Line
                  type="monotone"
                  dataKey="noDrops"
                  stroke="#0088fe"
                  name="No Drops"
                  dot={{ r: 3 }}
                  strokeWidth={2}
                  connectNulls
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Pi to Rho Section */}
        <div className="bg-white p-4 rounded shadow">
          <h2 className="text-xl font-semibold mb-6 text-center">
            Pi to Rho Throughput
          </h2>

          {/* Average Comparison */}
          <div className="mb-8">
            <h3 className="text-lg font-medium mb-3">Average Throughput</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={data.piRho.averages}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="name" />
                <YAxis
                  label={{ value: "MB/s", angle: -90, position: "insideLeft" }}
                />
                <Tooltip formatter={(value) => value.toFixed(2) + " MB/s"} />
                <Legend />
                <Bar dataKey="value" fill="#82ca9d" name="Throughput (MB/s)" />
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Line Chart Comparison */}
          <div>
            <h3 className="text-lg font-medium mb-3">Throughput per Sample</h3>
            <ResponsiveContainer width="100%" height={400}>
              <LineChart data={data.piRho.lineData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis
                  dataKey="index"
                  label={{
                    value: "Sample #",
                    position: "insideBottom",
                    offset: -5,
                  }}
                />
                <YAxis
                  label={{ value: "MB/s", angle: -90, position: "insideLeft" }}
                />
                <Tooltip
                  formatter={(value) =>
                    value !== null ? value.toFixed(2) + " MB/s" : "N/A"
                  }
                />
                <Legend />
                <Line
                  type="monotone"
                  dataKey="withDrops"
                  stroke="#ff7300"
                  name="With Drops"
                  dot={{ r: 3 }}
                  strokeWidth={2}
                  connectNulls
                />
                <Line
                  type="monotone"
                  dataKey="noDrops"
                  stroke="#0088fe"
                  name="No Drops"
                  dot={{ r: 3 }}
                  strokeWidth={2}
                  connectNulls
                />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>
    </div>
  );
}
