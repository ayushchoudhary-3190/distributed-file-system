'use client';

import { 
  HardDrive, 
  Users, 
  Clock, 
  Star, 
  Trash2, 
  Cloud,
  Settings,
  HelpCircle,
  Plus
} from 'lucide-react';
import { useState } from 'react';

interface SidebarProps {
  activeSection: string;
  onSectionChange: (section: string) => void;
}

const sections = [
  { id: 'my-drive', icon: HardDrive, label: 'My Drive' },
  { id: 'shared', icon: Users, label: 'Shared with me' },
  { id: 'recent', icon: Clock, label: 'Recent' },
  { id: 'starred', icon: Star, label: 'Starred' },
  { id: 'trash', icon: Trash2, label: 'Trash' },
];

export default function Sidebar({ activeSection, onSectionChange }: SidebarProps) {
  const [storageUsed] = useState(45);
  const storageGB = 15;
  const totalGB = 32;

  return (
    <div className="w-64 h-screen bg-[#1f1f1f] text-white flex flex-col flex-shrink-0">
      <div className="p-4">
        <div className="flex items-center gap-3 mb-6">
          <div className="w-10 h-10 bg-gradient-to-br from-green-400 to-blue-500 rounded-xl flex items-center justify-center">
            <Cloud className="w-6 h-6 text-white" />
          </div>
          <span className="text-xl font-medium">Drive</span>
        </div>

        <button className="w-full flex items-center gap-3 px-4 py-3 bg-white/10 hover:bg-white/20 rounded-full transition-colors mb-6">
          <Plus className="w-5 h-5" />
          <span className="font-medium">New</span>
        </button>
      </div>

      <nav className="flex-1 px-3">
        {sections.map((section) => (
          <button
            key={section.id}
            onClick={() => onSectionChange(section.id)}
            className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors mb-1 ${
              activeSection === section.id
                ? 'bg-blue-600 text-white'
                : 'text-gray-300 hover:bg-white/10'
            }`}
          >
            <section.icon className="w-5 h-5" />
            <span className="text-sm">{section.label}</span>
          </button>
        ))}
      </nav>

      <div className="p-4 border-t border-white/10">
        <button className="w-full flex items-center gap-3 px-3 py-2 text-gray-300 hover:bg-white/10 rounded-lg transition-colors mb-2">
          <Settings className="w-5 h-5" />
          <span className="text-sm">Settings</span>
        </button>
        <button className="w-full flex items-center gap-3 px-3 py-2 text-gray-300 hover:bg-white/10 rounded-lg transition-colors">
          <HelpCircle className="w-5 h-5" />
          <span className="text-sm">Help & Feedback</span>
        </button>
      </div>

      <div className="p-4 border-t border-white/10">
        <div className="flex items-center gap-2 mb-2">
          <Cloud className="w-4 h-4 text-gray-400" />
          <span className="text-xs text-gray-400">{storageGB} GB of {totalGB} GB used</span>
        </div>
        <div className="w-full h-2 bg-gray-700 rounded-full overflow-hidden">
          <div 
            className="h-full bg-gradient-to-r from-green-400 to-blue-500 rounded-full transition-all"
            style={{ width: `${storageUsed}%` }}
          />
        </div>
      </div>
    </div>
  );
}
