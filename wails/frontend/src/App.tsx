import {
  createEffect,
  createSignal,
  For,
  Show,
  type Component,
} from "solid-js";
import { Button } from "./components/ui/button";
import {
  createTreeCollection,
  Splitter,
  TreeCollection,
  TreeView,
} from "@ark-ui/solid";
import CheckSquareIcon from "lucide-solid/icons/check-square";
import ChevronRightIcon from "lucide-solid/icons/chevron-right";
import FileIcon from "lucide-solid/icons/file";
import FolderIcon from "lucide-solid/icons/folder";
import { ReadFolder } from "../wailsjs/go/main/App";
import { EventsOn } from "../wailsjs/runtime/runtime";
import { GetAppConfig } from "../wailsjs/go/main/App";
interface Node {
  Path: string;
  Name: string;
  Children?: Node[];
  IsDir: boolean;
}

// const collection = createTreeCollection<Node>({
//   nodeToValue: (node) => node.id,
//   nodeToString: (node) => node.name,
//   rootNode: {
//     id: "ROOT",
//     name: "",
//     children: [
//       {
//         id: "node_modules",
//         name: "node_modules",
//         children: [
//           { id: "node_modules/zag-js", name: "zag-js" },
//           { id: "node_modules/pandacss", name: "panda" },
//           {
//             id: "node_modules/@types",
//             name: "@types",
//             children: [
//               { id: "node_modules/@types/react", name: "react" },
//               { id: "node_modules/@types/react-dom", name: "react-dom" },
//             ],
//           },
//         ],
//       },
//       {
//         id: "src",
//         name: "src",
//         children: [
//           { id: "src/app.tsx", name: "app.tsx" },
//           { id: "src/index.ts", name: "index.ts" },
//         ],
//       },
//       { id: "panda.config", name: "panda.config.ts" },
//       { id: "package.json", name: "package.json" },
//       { id: "renovate.json", name: "renovate.json" },
//       { id: "readme.md", name: "README.md" },
//     ],
//   },
// });

const App: Component = () => {
  const createFolderStructure = async () => {
    const appConfig = await GetAppConfig();
    const rootNode = await ReadFolder(appConfig.Repository.Path);
    setCollection(
      createTreeCollection<Node>({
        nodeToValue: (node) => node.Path,
        nodeToString: (node) => node.Name,
        nodeToChildren: (node) => node?.Children!,
        rootNode,
      })
    );
  };
  createEffect(function changeFolderStructure() {
    mainBranch();
    createFolderStructure();
  });
  const [collection, setCollection] = createSignal<TreeCollection<Node>>(
    createTreeCollection<Node>({
      nodeToValue: (node) => node.Path,
      nodeToString: (node) => node.Name,
      nodeToChildren: (node) => node?.Children!,
      rootNode: {
        Path: "ROOT",
        Name: "",
        IsDir: true,
      },
    })
  );

  const [mainBranch, setMainBranch] = createSignal<string>("");
  EventsOn("branch-checkout", (branch: string) => {
    setMainBranch(branch);
  });

  return (
    <Splitter.Root
      defaultSize={[20, 80]}
      panels={[
        { id: "a", minSize: 10 },
        { id: "b", minSize: 30 },
      ]}
    >
      <Splitter.Context>
        {(api) => (
          <>
            <Splitter.Panel id="a">
              <div class="h-screen w-full bg-gray-50">
                <TreeView.Root collection={collection()}>
                  <TreeView.Label>Branch view</TreeView.Label>
                  <TreeView.Tree>
                    <For each={collection()?.rootNode?.Children}>
                      {(node, index) => (
                        <TreeNode node={node} indexPath={[index()]} />
                      )}
                    </For>
                  </TreeView.Tree>
                </TreeView.Root>
              </div>
            </Splitter.Panel>
            <Splitter.ResizeTrigger
              class="border border-gray-900"
              id="a:b"
              aria-label="Resize"
            />
            <Splitter.Panel id="b">
              <div class="h-screen w-full bg-gray-50">
                <Button
                  onClick={async () => {
                    const folder = await ReadFolder(
                      "/home/rnm/Dev/versionctrls-cli/watch-dir"
                    );
                    createFolderStructure();
                    // console.log(folder);
                  }}
                >
                  Read folder
                </Button>
              </div>
            </Splitter.Panel>
          </>
        )}
      </Splitter.Context>
    </Splitter.Root>
  );
};

const TreeNode = (props: TreeView.NodeProviderProps<Node>) => {
  const { node, indexPath } = props;
  return (
    <TreeView.NodeProvider node={node} indexPath={indexPath}>
      <Show
        when={node.Children && node.IsDir}
        fallback={
          <TreeView.Item>
            <div class="flex flex-row space-x-2">
              <TreeView.ItemIndicator>
                <CheckSquareIcon />
              </TreeView.ItemIndicator>
              <TreeView.ItemText>
                <div class="flex flex-row space-x-2">
                  <FileIcon class="mr-1" />
                  {node.Name}
                </div>
              </TreeView.ItemText>
            </div>
          </TreeView.Item>
        }
      >
        <TreeView.Branch>
          <TreeView.BranchControl>
            <div class="flex flex-row w-full space-x-2">
              <TreeView.BranchIndicator>
                <ChevronRightIcon />
              </TreeView.BranchIndicator>
              <TreeView.BranchText>
                <div class="flex flex-row">
                  <FolderIcon class="mr-1" /> {node.Name}
                </div>
              </TreeView.BranchText>
            </div>
          </TreeView.BranchControl>
          <TreeView.BranchContent>
            <TreeView.BranchIndentGuide />
            <For each={node.Children}>
              {(child, index) => (
                <TreeNode node={child} indexPath={[...indexPath, index()]} />
              )}
            </For>
          </TreeView.BranchContent>
        </TreeView.Branch>
      </Show>
    </TreeView.NodeProvider>
  );
};

export default App;
