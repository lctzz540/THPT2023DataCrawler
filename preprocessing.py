import pandas as pd

df = pd.read_csv("data2023.csv")
zero_columns = df.iloc[:, 1:].eq(0).all(axis=1)
df.drop_duplicates(subset=df.columns[0], keep="first", inplace=True)
filtered_df = df[~zero_columns]
filtered_df.to_csv("filtered_file.csv", index=False)
print(len(filtered_df))
