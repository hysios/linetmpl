import { Component } from "react";
import "./InlineEdit.css";

export interface ITree {
  name?: string;
  nodes?: INode[];
}

export interface INode {
  nodeType: string;
  pos: number;
  text?: string;
  args?: INode[];
  cmds?: INode[];
  ident?: string[];
  pipe?: INode;
}

export interface IInlineEditProps {
  text: string;
  preview?: boolean;
  tree?: ITree;
  context?: object;
  fields?: string[];
}

export default class InlineEditor extends Component<IInlineEditProps> {
  buildNode(node: INode, index: number) {
    switch (node.nodeType) {
      case "TEXT": // Text
        return this.textNode(node, index);
      case "ACTION": // ACTION
        return this.actionNode(node, index);
    }
  }

  textNode(node: INode, index: number) {
    return (
      <div className="flex-text" key={"node" + index}>
        <input
          type="text"
          size={node.text?.length}
          value={node.text}
          onChange={(e) => {
            node.text = e.target.value;
            this.setState({});
          }}
        ></input>
      </div>
    );
  }

  actionNode(node: INode, index: number) {
    if (node.pipe) {
      return this.pipeNode(node.pipe, index);
    }
  }

  pipeNode(node: INode, index: number) {
    return node.cmds?.map((cmd) => {
      console.log(cmd);
      switch (cmd?.nodeType) {
        case "COMMAND":
          return this.commandNode(cmd, index);
        default:
          return null;
      }
    });
  }

  commandNode(node: INode, index: number) {
    if (node.args && node.args.length > 0) {
      let first = node.args[0];
      console.log("first", first);
      switch (first.nodeType) {
        case "FIELD":
          return this.fieldNode(first, index);
      }
    }
    return null;
  }

  fieldNode(node: INode, index: number) {
    let { fields } = this.props,
      ident = node?.ident![0];

    return (
      <div className="flex-node" key={"node" + index}>
        <div className="node-box">
          <a className="node-close">×</a>
          <select className="btn-arrow-right">
            {fields?.map((field) => {
              return (
                <option selected={field == ident} title={field} value={field}>
                  {field}
                </option>
              );
            })}
          </select>
        </div>
        {/* <a className="node-btn btn-arrow-right">D</a> */}
        {/* <a className="node-btn pipe-icon">D</a>
        <a className="node-btn pipe-icon">D</a> */}
      </div>
    );
  }

  render() {
    let { text, tree } = this.props;

    return (
      <div>
        <div className="inline-editor flex-inline">
          {tree?.nodes?.map((node, index) => this.buildNode(node, index))}
        </div>
        {text}
      </div>
    );
  }
}

interface ITag {
  field: string;
  tip?: string;
  display?: string;
  selects?: string[];
  pipes?: IPipe[];
}

interface IPipe {}

class FieldTag extends Component<ITag> {
  render() {
    return <div></div>;
  }
}
