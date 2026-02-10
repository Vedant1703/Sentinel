"use client";

import { useState, useRef, useEffect } from "react";
import { Play, Square, ShieldAlert, ShieldCheck, Zap } from "lucide-react";
import { motion, AnimatePresence } from "framer-motion";

interface Log {
  id: number;
  status: number;
  time: string;
  msg: string;
}

export default function Playground() {
  const [logs, setLogs] = useState<Log[]>([]);
  const [isAttacking, setIsAttacking] = useState(false);
  const [stats, setStats] = useState({ sent: 0, allowed: 0, blocked: 0 });
  const attackRef = useRef<NodeJS.Timeout | null>(null);
  const scrollRef = useRef<HTMLDivElement>(null);

  // Auto-scroll logs
  useEffect(() => {
    if (scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [logs]);

  const addLog = (status: number) => {
    const now = new Date().toLocaleTimeString();
    const msg = status === 200 ? "Request Allowed" : "Rate Limit Exceeded";
    
    setLogs((prev) => {
      const newLog = { id: Date.now() + Math.random(), status, time: now, msg };
      return [...prev.slice(-49), newLog]; // Keep last 50 logs
    });

    setStats(prev => ({
      sent: prev.sent + 1,
      allowed: status === 200 ? prev.allowed + 1 : prev.allowed,
      blocked: status === 429 ? prev.blocked + 1 : prev.blocked
    }));
  };

  const fireRequest = async () => {
    try {
      // Use a specific path so we can rate limit it specifically
      const res = await fetch("http://localhost:8080/playground");
      if (res.ok) {
        addLog(200);
      } else {
        addLog(res.status);
      }
    } catch (error) {
      console.error(error);
      addLog(500);
    }
  };

  const toggleAttack = () => {
    if (isAttacking) {
      if (attackRef.current) clearInterval(attackRef.current);
      setIsAttacking(false);
    } else {
      setIsAttacking(true);
      // Fire a request every 50ms (20 RPS) -> Should trigger limit quickly
      attackRef.current = setInterval(fireRequest, 50);
    }
  };

  return (
    <div className="min-h-screen bg-background p-8 font-sans">
      <div className="mx-auto max-w-5xl">
        <header className="mb-12 text-center">
          <h1 className="text-4xl font-bold tracking-tight text-white mb-2">
            Traffic <span className="text-danger">Playground</span>
          </h1>
          <p className="text-gray-400">
            Simulate high-concurrency traffic to test the rate limiter.
          </p>
        </header>

        <div className="grid gap-8 md:grid-cols-2">
          {/* Controls Section */}
          <div className="space-y-6">
            <div className="rounded-xl border border-border bg-surface p-8 shadow-2xl text-center">
              <div className="mb-8 flex justify-center">
                <div className={`relative flex h-32 w-32 items-center justify-center rounded-full border-4 ${isAttacking ? "border-danger animate-pulse bg-danger/10" : "border-gray-700 bg-surface"}`}>
                  <Zap className={`h-16 w-16 ${isAttacking ? "text-danger" : "text-gray-600"}`} />
                </div>
              </div>

              <button
                onClick={toggleAttack}
                className={`group relative flex w-full items-center justify-center gap-3 rounded-lg px-8 py-4 text-lg font-bold transition-all ${
                  isAttacking
                    ? "bg-danger hover:bg-danger/90 text-white shadow-[0_0_20px_rgba(255,0,0,0.5)]"
                    : "bg-primary hover:bg-primary/90 text-black shadow-[0_0_20px_rgba(57,255,20,0.3)]"
                }`}
              >
                {isAttacking ? (
                  <>
                    <Square className="h-6 w-6 fill-current" /> STOP ATTACK
                  </>
                ) : (
                  <>
                    <Play className="h-6 w-6 fill-current" /> START SIMULATION
                  </>
                )}
              </button>
              <p className="mt-4 text-sm text-gray-500">
                Sends 20 requests/second to localhost:8080
              </p>
            </div>

            {/* Stats */}
            <div className="grid grid-cols-3 gap-4">
              <StatBox label="Sent" value={stats.sent} color="text-white" />
              <StatBox label="Allowed" value={stats.allowed} color="text-primary" />
              <StatBox label="Blocked" value={stats.blocked} color="text-danger" />
            </div>
          </div>

          {/* Terminal / Logs */}
          <div className="rounded-xl border border-border bg-[#0a0a0a] p-1 shadow-2xl overflow-hidden flex flex-col h-[500px]">
             <div className="flex items-center gap-2 border-b border-white/5 bg-white/5 px-4 py-3">
                <div className="flex gap-1.5">
                  <div className="h-3 w-3 rounded-full bg-red-500/50" />
                  <div className="h-3 w-3 rounded-full bg-yellow-500/50" />
                  <div className="h-3 w-3 rounded-full bg-green-500/50" />
                </div>
                <span className="ml-2 text-xs font-mono text-gray-500">live_traffic_logs</span>
             </div>
             
             <div 
               ref={scrollRef}
               className="flex-1 overflow-y-auto p-4 font-mono text-xs space-y-1 scrollbar-hide"
             >
               <AnimatePresence initial={false}>
                 {logs.map((log) => (
                   <motion.div
                      key={log.id}
                      initial={{ opacity: 0, x: -10 }}
                      animate={{ opacity: 1, x: 0 }}
                      className={`flex items-center gap-3 ${log.status === 200 ? "text-primary/80" : "text-danger/80"}`}
                   >
                      <span className="text-gray-600">[{log.time}]</span>
                      <span className={`font-bold ${log.status === 200 ? "bg-primary/20 px-1 rounded text-primary" : "bg-danger/20 px-1 rounded text-danger"}`}>
                        {log.status}
                      </span>
                      <span>{log.msg}</span>
                   </motion.div>
                 ))}
                 {logs.length === 0 && (
                   <div className="text-gray-600 italic">Waiting for traffic...</div>
                 )}
               </AnimatePresence>
             </div>
          </div>
        </div>
      </div>
    </div>
  );
}

function StatBox({ label, value, color }: { label: string, value: number, color: string }) {
  return (
    <div className="rounded-lg border border-border bg-surface p-4 text-center">
      <div className="text-xs text-gray-500 uppercase tracking-wider">{label}</div>
      <div className={`text-2xl font-bold ${color}`}>{value}</div>
    </div>
  )
}
