
import React, { useEffect, useRef } from 'react';
import { Terminal } from 'xterm';
import 'xterm/css/xterm.css';

const TerminalComponent: React.FC = () => {
  const terminalRef = useRef<HTMLDivElement | null>(null); 
  const terminalInstance = useRef<Terminal | null>(null); 

  useEffect(() => {
    const term = new Terminal(); 
    terminalInstance.current = term;

   
    if (terminalRef.current) {
      term.open(terminalRef.current);
    }

   
    const socket = new WebSocket(`ws://localhost:3000/ws`);

    
    socket.onmessage = (event: MessageEvent) => {
      term.write(event.data);
    };

    
    term.onData((data: string) => {
      socket.send(data);
    });

  
    return () => {
      socket.close();
      term.dispose();
    };
  }, []);

  return (
    <div
      ref={terminalRef}
      style={{
        width: '100%',
        height: '100vh',
        backgroundColor: 'black'
      }}
    ></div>
  );
};

export default TerminalComponent;
