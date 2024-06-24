import React, { useLayoutEffect } from 'react';
import * as am5 from "@amcharts/amcharts5";
import * as am5hierarchy from "@amcharts/amcharts5/hierarchy";
import am5themes_Animated from "@amcharts/amcharts5/themes/Animated";

// Generate and set data
// https://www.amcharts.com/docs/v5/charts/hierarchy/#Setting_data
// var maxLevels = 3;
// var maxNodes = 3;
// var maxValue = 100;

// var data = {
//     name: "Root",
//     children: []
// }
// generateLevel(data, "", 0);

// function generateLevel(data, name, level) {
//     for (var i = 0; i < Math.ceil(maxNodes * Math.random()) + 1; i++) {
//         var nodeName = name + "ABCDEFGHIJKLMNOPQRSTUVWXYZ"[i];
//         var child;
//         if (level < maxLevels) {
//             child = {
//                 name: nodeName + level
//             }

//             if (level > 0 && Math.random() < 0.5) {
//                 child.value = Math.round(Math.random() * maxValue);
//             }
//             else {
//                 child.children = [];
//                 generateLevel(child, nodeName + i, level + 1)
//             }
//         }
//         else {
//             child = {
//                 name: name + i,
//                 value: Math.round(Math.random() * maxValue)
//             }
//         }
//         data.children.push(child);
//     }

//     level++;
//     return data;
// }

// console.log(data)

const animalFind = (animals = {}, id, depth = 0, depthMax = 5) => {
    const animal = animals[id]
    if (!animal) {
        return { name: "unknown", gender: null, color: "#acacac", children: depth < depthMax ? [animalFind({}, null, depth+1), animalFind({}, null, depth+1)] : [] }
    }
    return { name: animal.name, gender: animal.gender, color:animal.gender == 1 ? "#e58f8f" : "#7bb4f1", children: depth < depthMax ? [
        animal?.parents?.[0] ? animalFind(animals, animal?.parents?.[0], depth+1) : animalFind({}, null, depth+1),
        animal?.parents?.[1] ? animalFind(animals, animal?.parents?.[1], depth+1) : animalFind({}, null, depth+1),
    ] : [] }
}

const AnimalTree = ({ animal = null }) => {
    console.log("animal", animal)
    useLayoutEffect(() => {
        let root = am5.Root.new("chartdiv");

        const data = { name: animal.name, children: [
            animal?.parents?.[0] ? animalFind(animal.family, animal?.parents?.[0], 1) : animalFind({}, null, 1),
            animal?.parents?.[1] ? animalFind(animal.family, animal?.parents?.[1], 1) : animalFind({}, null, 1),
        ] }

        console.log("chart data", data)

        root.setThemes([
            am5themes_Animated.new(root)
        ]);

        // Chart data

        var zoomableContainer = root.container.children.push(
            am5.ZoomableContainer.new(root, {
                width: am5.p100,
                height: am5.p100,
                wheelable: true,
                pinchZoom: true
            })
        );

        var zoomTools = zoomableContainer.children.push(am5.ZoomTools.new(root, {
            target: zoomableContainer
        }));

        // Create series
        // https://www.amcharts.com/docs/v5/charts/hierarchy/#Adding
        var series = zoomableContainer.contents.children.push(am5hierarchy.Tree.new(root, {
            singleBranchOnly: false,
            downDepth: 1,
            initialDepth: 10,
            valueField: "value",
            categoryField: "name",
            childDataField: "children",
            fillField: "color"
        }));

        series.labels.template.set("minScale", 0);

        series.data.setAll([data]);
        series.set("selectedDataItem", series.dataItems[0]);

        // Chart stop

        return () => {
            root.dispose();
        };
    }, []);

    return (
        <div id="chartdiv" style={{ width: "100%", height: "100%" }}></div>
    );
}

export default AnimalTree