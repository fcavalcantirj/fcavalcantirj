import json
import struct

def get_mesh_names_from_glb(glb_file_path):
    """Extracts mesh names from a GLB file."""
    try:
        with open(glb_file_path, "rb") as f:
            # Read GLB header
            magic = f.read(4)
            if magic != b'glTF':
                print("Error: Not a valid GLB file (magic number incorrect).")
                return []

            version = struct.unpack('<I', f.read(4))[0]
            if version != 2:
                print(f"Warning: GLB version is {version}, expected 2. Attempting to parse anyway.")

            length = struct.unpack('<I', f.read(4))[0]

            # Read JSON chunk header
            json_chunk_length = struct.unpack('<I', f.read(4))[0]
            json_chunk_type = f.read(4)

            if json_chunk_type != b'JSON':
                print("Error: First chunk is not JSON.")
                return []

            # Read JSON content
            json_data_str = f.read(json_chunk_length).decode('utf-8')
            gltf_json = json.loads(json_data_str)

            mesh_names = []
            if 'meshes' in gltf_json:
                for i, mesh in enumerate(gltf_json['meshes']):
                    name = mesh.get('name', f'Unnamed_Mesh_{i}')
                    mesh_names.append(name)
                    print(f"Found mesh: Name: 	'{name}'")
            
            # Attempt to find nodes that reference meshes to get more context if names are generic
            if 'nodes' in gltf_json:
                print("\nNodes referencing meshes (if any):")
                for i, node in enumerate(gltf_json['nodes']):
                    node_name = node.get('name', f'Unnamed_Node_{i}')
                    if 'mesh' in node: # node.mesh is the index of the mesh in gltf_json['meshes']
                        mesh_index = node['mesh']
                        if mesh_index < len(gltf_json.get('meshes', [])):
                            referenced_mesh_name = gltf_json['meshes'][mesh_index].get('name', f'Unnamed_Mesh_{mesh_index}')
                            print(f"Node: 	'{node_name}' 	(references mesh: '{referenced_mesh_name}' at index {mesh_index})")
                        else:
                            print(f"Node: 	'{node_name}' 	(references out-of-bounds mesh index {mesh_index})")
                    elif node_name and any(keyword in node_name.lower() for keyword in ['bus', 'window', 'car', 'building', 'door']):
                         print(f"Node: 	'{node_name}' 	(does not directly reference a mesh but name is interesting)")


            if not mesh_names:
                print("No meshes found in the JSON part of the GLB.")
            return mesh_names

    except FileNotFoundError:
        print(f"Error: File not found at {glb_file_path}")
        return []
    except Exception as e:
        print(f"An error occurred: {e}")
        return []

if __name__ == "__main__":
    model_file = "/home/ubuntu/threejs_clone/models/LittlestTokyo.glb"
    print(f"Inspecting GLB file: {model_file}\n")
    names = get_mesh_names_from_glb(model_file)
    if names:
        print("\nSuccessfully extracted mesh names.")
    else:
        print("\nFailed to extract mesh names or no meshes found.")

