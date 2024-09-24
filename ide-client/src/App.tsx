import React from 'react';

import TerminalComponent from './components/terminal';
import FileTree from './components/fileTree';
import { FileTreeProvider } from './Context/FileTreeContext';

const App: React.FC = () => {
  return (
    <FileTreeProvider>
    <div className="App">
      
      <FileTree/>
    </div>
    </FileTreeProvider>
  );
};

export default App;
