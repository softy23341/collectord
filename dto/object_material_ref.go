package dto

//go:generate dbgen -type ObjectMaterialRef

// ObjectMaterialRef TBD
type ObjectMaterialRef struct {
	ObjectID   int64 `db:"object_id"`
	MaterialID int64 `db:"material_id"`
}

// ObjectMaterialRefList TBD
type ObjectMaterialRefList []*ObjectMaterialRef

// ObjectIDToMaterialListIDsMap TBD
func (ol ObjectMaterialRefList) ObjectIDToMaterialListIDsMap() map[int64][]int64 {
	objectIDToMaterials := make(map[int64][]int64, 0)
	for _, ref := range ol {
		objectIDToMaterials[ref.ObjectID] = append(objectIDToMaterials[ref.ObjectID], ref.MaterialID)
	}
	return objectIDToMaterials
}

// MaterialsIDs TBD
func (ol ObjectMaterialRefList) MaterialsIDs() []int64 {
	materialsIDsMap := make(map[int64]struct{}, len(ol))
	for _, materialRef := range ol {
		materialsIDsMap[materialRef.MaterialID] = struct{}{}
	}

	var materialsIDsList []int64
	for actorID := range materialsIDsMap {
		materialsIDsList = append(materialsIDsList, actorID)
	}

	return materialsIDsList
}
