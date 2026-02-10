"use client";

import { useEffect, useState } from "react";
import { Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { Activity, ShieldAlert, ShieldCheck } from "lucide-react";
import { motion } from "framer-motion";

// Types for our metrics
interface Metrics {
  allowed_requests: number;
  blocked_requests: number;
  redis_errors: number;
}

interface DataPoint {
  time: string;
  allowed: number;
  blocked: number;
}

export default function Home() {
  const [metrics, setMetrics] = useState<Metrics>({ allowed_requests: 0, blocked_requests: 0, redis_errors: 0 });
  const [history, setHistory] = useState<DataPoint[]>([]);
  const [isConnected, setIsConnected] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const res = await fetch("http://localhost:8080/metrics");
        if (res.ok) {
          const data: Metrics = await res.json();
          setMetrics(data);
          setIsConnected(true);

          setHistory((prev) => {
            const now = new Date().toLocaleTimeString();
            const newPoint = { 
              time: now, 
              allowed: data.allowed_requests, 
              blocked: data.blocked_requests 
            };
            // Keep last 20 points
            const newHistory = [...prev, newPoint];
            return newHistory.slice(-20);
          });
        }
      } catch (error) {
        console.error("Failed to fetch metrics", error);
        setIsConnected(false);
      }
    };

    const interval = setInterval(fetchData, 1000);
    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-background p-8 font-sans">
      <header className="mb-8 flex items-center justify-between">
        <div>
          <h1 className="text-4xl font-bold tracking-tight text-white">Sentinel <span className="text-primary">Dashboard</span></h1>
          <p className="text-gray-400">Distributed Rate Limiting System</p>
        </div>
        <div className={`flex items-center gap-2 rounded-full px-4 py-1 text-sm font-medium ${isConnected ? "bg-primary/20 text-primary" : "bg-danger/20 text-danger"}`}>
          <div className={`h-2 w-2 rounded-full ${isConnected ? "bg-primary animate-pulse" : "bg-danger"}`} />
          {isConnected ? "System Online" : "System Offline"}
        </div>
      </header>

      {/* Metrics Grid */}
      <div className="grid gap-6 md:grid-cols-3 mb-8">
        <MetricCard 
          title="Allowed Requests" 
          value={metrics.allowed_requests} 
          icon={<ShieldCheck className="h-6 w-6 text-primary" />} 
          color="text-primary"
        />
        <MetricCard 
          title="Blocked Requests" 
          value={metrics.blocked_requests} 
          icon={<ShieldAlert className="h-6 w-6 text-danger" />} 
          color="text-danger"
        />
        <MetricCard 
          title="Redis Errors" 
          value={metrics.redis_errors} 
          icon={<Activity className="h-6 w-6 text-warning" />} 
          color="text-warning"
        />
      </div>

      {/* Main Chart */}
      <div className="rounded-xl border border-border bg-surface p-6 shadow-2xl">
        <h2 className="mb-6 text-xl font-semibold text-white">Real-Time Traffic Visualization</h2>
        <div className="h-[400px] w-full">
          <ResponsiveContainer width="100%" height="100%">
            <LineChart data={history}>
              <XAxis dataKey="time" stroke="#666" fontSize={12} tickLine={false} axisLine={false} />
              <YAxis stroke="#666" fontSize={12} tickLine={false} axisLine={false} />
              <Tooltip 
                contentStyle={{ backgroundColor: "#111", border: "1px solid #333" }}
                itemStyle={{ color: "#fff" }}
              />
              <Line 
                type="monotone" 
                dataKey="allowed" 
                stroke="#39FF14" 
                strokeWidth={2} 
                dot={false} 
                animationDuration={300}
              />
              <Line 
                type="monotone" 
                dataKey="blocked" 
                stroke="#FF0000" 
                strokeWidth={2} 
                dot={false} 
                animationDuration={300}
              />
            </LineChart>
          </ResponsiveContainer>
        </div>
      </div>
    </div>
  );
}

function MetricCard({ title, value, icon, color }: { title: string, value: number, icon: React.ReactNode, color: string }) {
  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      className="rounded-xl border border-border bg-surface p-6 shadow-lg"
    >
      <div className="flex items-center justify-between">
        <h3 className="text-sm font-medium text-gray-400">{title}</h3>
        {icon}
      </div>
      <div className={`mt-2 text-3xl font-bold ${color}`}>
        {value.toLocaleString()}
      </div>
    </motion.div>
  );
}
