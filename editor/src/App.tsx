import React, { useEffect, useState } from "react";
import logo from "./logo.svg";
import "./App.css";
import InlineEditor, { ITree } from "./InlineEdit";

function App() {
  var [root, setRoot] = useState<ITree>({});
  useEffect(() => {
    async function fetchMyAPI() {
      let data = (await fetch("/1").then((res) => res.json())) as any;
      setRoot(data["data"]);
      console.log(data["data"]);
    }

    fetchMyAPI();
  }, []);

  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <p>
          Edit <code>src/App.tsx</code> and save to reload.
        </p>
        <InlineEditor
          tree={root}
          fields={[
            "日期",
            "车牌",
            "年",
            "月",
            "日",
            "违法代码",
            "序号",
            "设备编号",
            "违法行为",
          ]}
          text={"/[.日期]/[.车牌]/[.年]/[.月]/[.日]/[.违法代码]_[.序号].jpg"}
        />
      </header>
    </div>
  );
}

export default App;
