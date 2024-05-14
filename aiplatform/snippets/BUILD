go_library(
    name = "snippets",
    srcs = [
        "create_dataset.go",
        "delete_dataset.go",
        "embeddings.go",
        "embeddings_preview.go",
        "import_data_image_classification.go",
        "list_models.go",
        "test.go",
    ],
)

go_test(
    name = "snippets_test",
    srcs = [
        "create_dataset_test.go",
        "delete_dataset_test.go",
        "embeddings_test.go",
        "import_data_image_classification_test.go",
        "list_models_test.go",
    ],
    library = ":snippets",
)
