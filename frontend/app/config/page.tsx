"use client";

import { useState } from "react";
import { Save, RefreshCw, CheckCircle2 } from "lucide-react";

export default function ConfigPage() {
  const [loading, setLoading] = useState(false);
  const [success, setSuccess] = useState<string | null>(null);
  
  // Default values matching backend defaults
  const [config, setConfig] = useState({
    path: "/playground",
    limit: 10,
    window: 60
  });

  const handleUpdate = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setSuccess(null);

    try {
      const res = await fetch(`${process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"}/api/config`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(config),
      });

      if (res.ok) {
        setSuccess(`Updated rule for ${config.path}`);
        setTimeout(() => setSuccess(null), 3000);
      } else {
        alert("Failed to update config");
      }
    } catch (error) {
      console.error(error);
      alert("Error connecting to server");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen bg-background p-8 font-sans">
      <div className="mx-auto max-w-2xl">
        <header className="mb-12 text-center">
          <h1 className="text-4xl font-bold tracking-tight text-white mb-2">
            System <span className="text-warning">Configuration</span>
          </h1>
          <p className="text-gray-400">
            Dynamically update rate limiting rules without restarting.
          </p>
        </header>

        <div className="rounded-xl border border-border bg-surface p-8 shadow-2xl">
          <form onSubmit={handleUpdate} className="space-y-6">
            
            {/* Route Selection */}
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">Target Route</label>
              <select 
                value={config.path}
                onChange={(e) => setConfig({ ...config, path: e.target.value })}
                className="w-full rounded-lg bg-black/50 border border-border px-4 py-3 text-white focus:border-primary focus:outline-none"
              >
                <option value="/playground">/playground (Test Route)</option>
                <option value="/login">/login</option>
                <option value="/search">/search</option>
                <option value="default">Default (All other routes)</option>
              </select>
            </div>

            {/* Limit Input */}
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">
                Request Limit (Requests / Window)
              </label>
              <input 
                type="number" 
                value={config.limit}
                onChange={(e) => setConfig({ ...config, limit: parseInt(e.target.value) })}
                className="w-full rounded-lg bg-black/50 border border-border px-4 py-3 text-white focus:border-primary focus:outline-none"
                min="1"
              />
            </div>

            {/* Window Input */}
            <div>
              <label className="block text-sm font-medium text-gray-400 mb-2">
                Time Window (Seconds)
              </label>
              <input 
                type="number" 
                value={config.window}
                onChange={(e) => setConfig({ ...config, window: parseInt(e.target.value) })}
                className="w-full rounded-lg bg-black/50 border border-border px-4 py-3 text-white focus:border-primary focus:outline-none"
                min="1"
              />
            </div>

            {/* Submit Button */}
            <button
              type="submit"
              disabled={loading}
              className={`flex w-full items-center justify-center gap-2 rounded-lg py-4 text-lg font-bold transition-all ${
                success 
                  ? "bg-primary text-black"
                  : "bg-white text-black hover:bg-gray-200"
              }`}
            >
              {loading ? (
                <RefreshCw className="h-6 w-6 animate-spin" />
              ) : success ? (
                <>
                  <CheckCircle2 className="h-6 w-6" /> Updated!
                </>
              ) : (
                <>
                  <Save className="h-6 w-6" /> Save Configuration
                </>
              )}
            </button>
          </form>
        </div>
        
        <div className="mt-8 rounded-lg border border-primary/20 bg-primary/5 p-4 text-sm text-primary/80">
          <strong>Tip:</strong> Set the limit to something low (e.g., 5 req / 10s) and test it in the Playground!
        </div>
      </div>
    </div>
  );
}
